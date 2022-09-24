先说结论：当发现存在和newRS Template一致的oldRS时，会更新oldRS的annotations为最新的版本号，并将旧的版本号追加annotations deployment.kubernetes.io/revision-history

某平台统计deployment历史版本号发现中间版本缺失

![image-20220924193238122](https://raw.githubusercontent.com/helen-frank/img/master/img/202209241932183.png)

环境上

```shell
# kubectl -l app=nginx475 get rs -oyaml | grep deployment.kubernetes.io/revision:
      deployment.kubernetes.io/revision: "8"
      deployment.kubernetes.io/revision: "5"
      deployment.kubernetes.io/revision: "1"
      deployment.kubernetes.io/revision: "7"
      deployment.kubernetes.io/revision: "9"
      deployment.kubernetes.io/revision: "10"
      deployment.kubernetes.io/revision: "4"
      deployment.kubernetes.io/revision: "6"
```

可以看到平台上的显示是没有问题的，那么为什么会发生这种情况呢？

k8s release1.18

`pkg/controller/deployment/deployment_controller.go`

```go
// NewDeploymentController creates a new DeploymentController.
func NewDeploymentController(dInformer appsinformers.DeploymentInformer, rsInformer appsinformers.ReplicaSetInformer, podInformer coreinformers.PodInformer, client clientset.Interface) (*DeploymentController, error) {
 ...
  dc.syncHandler = dc.syncDeployment
 ...v
}
```

`pkg/controller/deployment/deployment_controller.go`

```go
// syncDeployment will sync the deployment with the given key.
// This function is not meant to be invoked concurrently with the same key.
func (dc *DeploymentController) syncDeployment(key string) error {
	...
  	switch d.Spec.Strategy.Type {
	case apps.RecreateDeploymentStrategyType:
		return dc.rolloutRecreate(d, rsList, podMap)
	case apps.RollingUpdateDeploymentStrategyType:
		return dc.rolloutRolling(d, rsList)
	}
	...
}
```

`pkg/controller/deployment/recreate.go`

```go
// rolloutRecreate implements the logic for recreating a replica set.
func (dc *DeploymentController) rolloutRecreate(d *apps.Deployment, rsList []*apps.ReplicaSet, podMap map[types.UID][]*v1.Pod) error {
	// Don't create a new RS if not already existed, so that we avoid scaling up before scaling down.
	newRS, oldRSs, err := dc.getAllReplicaSetsAndSyncRevision(d, rsList, false)
	if err != nil {
		return err
	}
 ... 
}
```

`pkg/controller/deployment/sync.go`

```go
// getAllReplicaSetsAndSyncRevision returns all the replica sets for the provided deployment (new and all old), with new RS's and deployment's revision updated.
//
// rsList should come from getReplicaSetsForDeployment(d).
//
// 1. Get all old RSes this deployment targets, and calculate the max revision number among them (maxOldV).
// 2. Get new RS this deployment targets (whose pod template matches deployment's), and update new RS's revision number to (maxOldV + 1),
//    only if its revision number is smaller than (maxOldV + 1). If this step failed, we'll update it in the next deployment sync loop.
// 3. Copy new RS's revision number to deployment (update deployment's revision). If this step failed, we'll update it in the next deployment sync loop.
//
// Note that currently the deployment controller is using caches to avoid querying the server for reads.
// This may lead to stale reads of replica sets, thus incorrect deployment status.
func (dc *DeploymentController) getAllReplicaSetsAndSyncRevision(d *apps.Deployment, rsList []*apps.ReplicaSet, createIfNotExisted bool) (*apps.ReplicaSet, []*apps.ReplicaSet, error) {
	_, allOldRSs := deploymentutil.FindOldReplicaSets(d, rsList)

	// Get new replica set with the updated revision number
	newRS, err := dc.getNewReplicaSet(d, rsList, allOldRSs, createIfNotExisted)
	if err != nil {
		return nil, nil, err
	}

	return newRS, allOldRSs, nil
}
```

`pkg/controller/deployment/sync.go`

```go
// Returns a replica set that matches the intent of the given deployment. Returns nil if the new replica set doesn't exist yet.
// 1. Get existing new RS (the RS that the given deployment targets, whose pod template is the same as deployment's).
// 2. If there's existing new RS, update its revision number if it's smaller than (maxOldRevision + 1), where maxOldRevision is the max revision number among all old RSes.
// 3. If there's no existing new RS and createIfNotExisted is true, create one with appropriate revision number (maxOldRevision + 1) and replicas.
// Note that the pod-template-hash will be added to adopted RSes and pods.
func (dc *DeploymentController) getNewReplicaSet(d *apps.Deployment, rsList, oldRSs []*apps.ReplicaSet, createIfNotExisted bool) (*apps.ReplicaSet, error) {
	existingNewRS := deploymentutil.FindNewReplicaSet(d, rsList)

	// Calculate the max revision number among all old RSes
	maxOldRevision := deploymentutil.MaxRevision(oldRSs)
	// Calculate revision number for this new replica set
	newRevision := strconv.FormatInt(maxOldRevision+1, 10)

	// Latest replica set exists. We need to sync its annotations (includes copying all but
	// annotationsToSkip from the parent deployment, and update revision, desiredReplicas,
	// and maxReplicas) and also update the revision annotation in the deployment with the
	// latest revision.
	if existingNewRS != nil {
		rsCopy := existingNewRS.DeepCopy()

		// Set existing new replica set's annotation
		annotationsUpdated := deploymentutil.SetNewReplicaSetAnnotations(d, rsCopy, newRevision, true, maxRevHistoryLengthInChars)
		minReadySecondsNeedsUpdate := rsCopy.Spec.MinReadySeconds != d.Spec.MinReadySeconds
		if annotationsUpdated || minReadySecondsNeedsUpdate {
			rsCopy.Spec.MinReadySeconds = d.Spec.MinReadySeconds
			return dc.client.AppsV1().ReplicaSets(rsCopy.ObjectMeta.Namespace).Update(context.TODO(), rsCopy, metav1.UpdateOptions{})
		}

		// Should use the revision in existingNewRS's annotation, since it set by before
		needsUpdate := deploymentutil.SetDeploymentRevision(d, rsCopy.Annotations[deploymentutil.RevisionAnnotation])
		// If no other Progressing condition has been recorded and we need to estimate the progress
		// of this deployment then it is likely that old users started caring about progress. In that
		// case we need to take into account the first time we noticed their new replica set.
		cond := deploymentutil.GetDeploymentCondition(d.Status, apps.DeploymentProgressing)
		if deploymentutil.HasProgressDeadline(d) && cond == nil {
			msg := fmt.Sprintf("Found new replica set %q", rsCopy.Name)
			condition := deploymentutil.NewDeploymentCondition(apps.DeploymentProgressing, v1.ConditionTrue, deploymentutil.FoundNewRSReason, msg)
			deploymentutil.SetDeploymentCondition(&d.Status, *condition)
			needsUpdate = true
		}

		if needsUpdate {
			var err error
			if d, err = dc.client.AppsV1().Deployments(d.Namespace).UpdateStatus(context.TODO(), d, metav1.UpdateOptions{}); err != nil {
				return nil, err
			}
		}
		return rsCopy, nil
	}

	if !createIfNotExisted {
		return nil, nil
	}

	// new ReplicaSet does not exist, create one.
	newRSTemplate := *d.Spec.Template.DeepCopy()
	podTemplateSpecHash := controller.ComputeHash(&newRSTemplate, d.Status.CollisionCount)
	newRSTemplate.Labels = labelsutil.CloneAndAddLabel(d.Spec.Template.Labels, apps.DefaultDeploymentUniqueLabelKey, podTemplateSpecHash)
	// Add podTemplateHash label to selector.
	newRSSelector := labelsutil.CloneSelectorAndAddLabel(d.Spec.Selector, apps.DefaultDeploymentUniqueLabelKey, podTemplateSpecHash)

	// Create new ReplicaSet
	newRS := apps.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{
			// Make the name deterministic, to ensure idempotence
			Name:            d.Name + "-" + podTemplateSpecHash,
			Namespace:       d.Namespace,
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(d, controllerKind)},
			Labels:          newRSTemplate.Labels,
		},
		Spec: apps.ReplicaSetSpec{
			Replicas:        new(int32),
			MinReadySeconds: d.Spec.MinReadySeconds,
			Selector:        newRSSelector,
			Template:        newRSTemplate,
		},
	}
	allRSs := append(oldRSs, &newRS)
	newReplicasCount, err := deploymentutil.NewRSNewReplicas(d, allRSs, &newRS)
	if err != nil {
		return nil, err
	}

	*(newRS.Spec.Replicas) = newReplicasCount
	// Set new replica set's annotation
	deploymentutil.SetNewReplicaSetAnnotations(d, &newRS, newRevision, false, maxRevHistoryLengthInChars)
	// Create the new ReplicaSet. If it already exists, then we need to check for possible
	// hash collisions. If there is any other error, we need to report it in the status of
	// the Deployment.
	alreadyExists := false
	createdRS, err := dc.client.AppsV1().ReplicaSets(d.Namespace).Create(context.TODO(), &newRS, metav1.CreateOptions{})
	switch {
	// We may end up hitting this due to a slow cache or a fast resync of the Deployment.
	case errors.IsAlreadyExists(err):
		alreadyExists = true

		// Fetch a copy of the ReplicaSet.
		rs, rsErr := dc.rsLister.ReplicaSets(newRS.Namespace).Get(newRS.Name)
		if rsErr != nil {
			return nil, rsErr
		}

		// If the Deployment owns the ReplicaSet and the ReplicaSet's PodTemplateSpec is semantically
		// deep equal to the PodTemplateSpec of the Deployment, it's the Deployment's new ReplicaSet.
		// Otherwise, this is a hash collision and we need to increment the collisionCount field in
		// the status of the Deployment and requeue to try the creation in the next sync.
		controllerRef := metav1.GetControllerOf(rs)
		if controllerRef != nil && controllerRef.UID == d.UID && deploymentutil.EqualIgnoreHash(&d.Spec.Template, &rs.Spec.Template) {
			createdRS = rs
			err = nil
			break
		}

		// Matching ReplicaSet is not equal - increment the collisionCount in the DeploymentStatus
		// and requeue the Deployment.
		if d.Status.CollisionCount == nil {
			d.Status.CollisionCount = new(int32)
		}
		preCollisionCount := *d.Status.CollisionCount
		*d.Status.CollisionCount++
		// Update the collisionCount for the Deployment and let it requeue by returning the original
		// error.
		_, dErr := dc.client.AppsV1().Deployments(d.Namespace).UpdateStatus(context.TODO(), d, metav1.UpdateOptions{})
		if dErr == nil {
			klog.V(2).Infof("Found a hash collision for deployment %q - bumping collisionCount (%d->%d) to resolve it", d.Name, preCollisionCount, *d.Status.CollisionCount)
		}
		return nil, err
	case errors.HasStatusCause(err, v1.NamespaceTerminatingCause):
		// if the namespace is terminating, all subsequent creates will fail and we can safely do nothing
		return nil, err
	case err != nil:
		msg := fmt.Sprintf("Failed to create new replica set %q: %v", newRS.Name, err)
		if deploymentutil.HasProgressDeadline(d) {
			cond := deploymentutil.NewDeploymentCondition(apps.DeploymentProgressing, v1.ConditionFalse, deploymentutil.FailedRSCreateReason, msg)
			deploymentutil.SetDeploymentCondition(&d.Status, *cond)
			// We don't really care about this error at this point, since we have a bigger issue to report.
			// TODO: Identify which errors are permanent and switch DeploymentIsFailed to take into account
			// these reasons as well. Related issue: https://github.com/kubernetes/kubernetes/issues/18568
			_, _ = dc.client.AppsV1().Deployments(d.Namespace).UpdateStatus(context.TODO(), d, metav1.UpdateOptions{})
		}
		dc.eventRecorder.Eventf(d, v1.EventTypeWarning, deploymentutil.FailedRSCreateReason, msg)
		return nil, err
	}
	if !alreadyExists && newReplicasCount > 0 {
		dc.eventRecorder.Eventf(d, v1.EventTypeNormal, "ScalingReplicaSet", "Scaled up replica set %s to %d", createdRS.Name, newReplicasCount)
	}

	needsUpdate := deploymentutil.SetDeploymentRevision(d, newRevision)
	if !alreadyExists && deploymentutil.HasProgressDeadline(d) {
		msg := fmt.Sprintf("Created new replica set %q", createdRS.Name)
		condition := deploymentutil.NewDeploymentCondition(apps.DeploymentProgressing, v1.ConditionTrue, deploymentutil.NewReplicaSetReason, msg)
		deploymentutil.SetDeploymentCondition(&d.Status, *condition)
		needsUpdate = true
	}
	if needsUpdate {
		_, err = dc.client.AppsV1().Deployments(d.Namespace).UpdateStatus(context.TODO(), d, metav1.UpdateOptions{})
	}
	return createdRS, err
}
```

`pkg/controller/deployment/util/deployment_util.go`

```go
// FindNewReplicaSet returns the new RS this given deployment targets (the one with the same pod template).
func FindNewReplicaSet(deployment *apps.Deployment, rsList []*apps.ReplicaSet) *apps.ReplicaSet {
	sort.Sort(controller.ReplicaSetsByCreationTimestamp(rsList))
	for i := range rsList {
		if EqualIgnoreHash(&rsList[i].Spec.Template, &deployment.Spec.Template) {
			// In rare cases, such as after cluster upgrades, Deployment may end up with
			// having more than one new ReplicaSets that have the same template as its template,
			// see https://github.com/kubernetes/kubernetes/issues/40415
			// We deterministically choose the oldest new ReplicaSet.
			return rsList[i]
		}
	}
	// new ReplicaSet does not exist.
	return nil
}
```

由上代码可知：

当发现存在和newRS Template一致的oldRS时，会更新oldRS的annotations为最新的版本号，并将旧的版本号追加annotations deployment.kubernetes.io/revision-history （xxx,xxx,xxx,oldRSVersion）,长度限制为2000字符（release-1.18）,当超出2000之后会报Warning`Not appending revision due to length limit of 2000 reached`

因为是更新原有rs而来所以表现为部分版本号缺失