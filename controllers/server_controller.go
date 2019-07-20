/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	core "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	minecraftv1alpha1 "github.com/jbeda/kinecraft/api/v1alpha1"
)

// ServerReconciler reconciles a Server object
type ServerReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func ignoreNotFound(err error) error {
	if apierrs.IsNotFound(err) {
		return nil
	}
	return err
}

// +kubebuilder:rbac:groups=minecraft.tgik.io,resources=servers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=minecraft.tgik.io,resources=servers/status,verbs=get;update;patch
// TODO(jbeda): list RBAC stuff

func (r *ServerReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("server", req.NamespacedName)

	var mcServer minecraftv1alpha1.Server
	if err := r.Get(ctx, req.NamespacedName, &mcServer); err != nil {
		log.Error(err, "unable to fetch Server")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, ignoreNotFound(err)
	}

	// TODO(jbeda): List out the Pods that belong to this server and updatedate
	// the status field.  If we already have a server running then exit out here.

	pod, err := r.constructPod(&mcServer)
	if err != nil {
		return ctrl.Result{}, err
	}

	if err := r.Create(ctx, pod); err != nil {
		log.Error(err, "unable to create Pod for Server", "pod", pod)
		return ctrl.Result{}, err
	}

	log.V(1).Info("created Pod for Server run", "pod", pod)
	return ctrl.Result{}, nil
}

func (r *ServerReconciler) constructPod(s *minecraftv1alpha1.Server) (*core.Pod, error) {
	namePrefix := fmt.Sprintf("mc-%s", s.Name)
	pod := &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
			Name:        namePrefix,
			// TODO GeneratedName -- we can add this back when we test for existence
			Namespace: s.Namespace,
		},
		// containers:
		// - image: itzg/minecraft-server
		// 	name: minecraft-server
		// 	env:
		// 	- name: TYPE
		// 		value: "VANILLA"
		// 	- name: EULA
		// 		value: "TRUE"
		// 	ports:
		// 	- containerPort: 25565
		// 		name: minecraft
		// 		protocol: TCP
		Spec: core.PodSpec{
			Containers: []core.Container{
				core.Container{
					Image: "itzg/minecraft-server",
					Name:  "minecraft-server",
					Env:   []core.EnvVar{},
					Ports: []core.ContainerPort{
						core.ContainerPort{
							ContainerPort: 25565,
							Name:          "minecraft",
							Protocol:      "TCP",
						},
					},
				},
			},
		},
	}

	addEnv := func(key string, value string) {
		pod.Spec.Containers[0].Env = append(pod.Spec.Containers[0].Env,
			core.EnvVar{Name: key, Value: value})
	}
	bool2string := func(b bool) string {
		if b {
			return "TRUE"
		} else {
			return "FALSE"
		}
	}

	// TODO: If these values are blank we should just not set the env variable.
	addEnv("EULA", bool2string(s.Spec.EULA))
	addEnv("TYPE", s.Spec.ServerType)
	addEnv("SERVER_NAME", s.Spec.ServerName)
	addEnv("OPS", strings.Join(s.Spec.Ops, ","))
	addEnv("WHITELIST", strings.Join(s.Spec.Allowlist, ","))

	if err := ctrl.SetControllerReference(pod, s, r.Scheme); err != nil {
		return nil, err
	}
	return pod, nil
}

func (r *ServerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// TODO(jbeda): Make sure that we are getting kicked on pod changes also.
	return ctrl.NewControllerManagedBy(mgr).
		For(&minecraftv1alpha1.Server{}).
		Complete(r)
}
