package main

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	typedv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	v1 "k8s.io/client-go/listers/core/v1"
)

const (
	coreV1APIVersion = "v1"

	namespaceKind = "Namespace"
	podKind       = "Pod"

	appNamespace = "app"
	appName      = "hello-world"

	creationWaitTimeout = 10 * time.Second
)

type program struct {
	namespaceClient typedv1.NamespaceInterface
	podClient       typedv1.PodInterface

	namespaceLister v1.NamespaceLister
	podLister       v1.PodLister
}

func newProgram(clientSet kubernetes.Interface, namespaceLister v1.NamespaceLister, podLister v1.PodLister) *program {
	return &program{
		namespaceClient: clientSet.CoreV1().Namespaces(),
		podClient:       clientSet.CoreV1().Pods(appNamespace),
		namespaceLister: namespaceLister,
		podLister:       podLister,
	}
}

func (p *program) run(ctx context.Context, delete bool) {
	if delete {
		fmt.Println("6. delete the hello-world pod")
		if err := p.deleteHelloWorldPod(ctx); err != nil {
			fmt.Println(err)
		}
		return
	}

	fmt.Println("2. print out the namespaces in the cluster")
	if err := p.listNamespaces(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("3. create a new namespace")
	if err := p.createAppNamespace(ctx); err != nil {
		fmt.Println(err)
	}
	fmt.Println("4. create a pod in the new namespace that runs a simple hello-world container")
	if err := p.createHelloWorldPod(ctx); err != nil {
		fmt.Println(err)
	}
	fmt.Println("5. print out the pod names and namespaces for pods with the label ‘k8s-app=kube-dns’")
	if err := p.listDNSPods(); err != nil {
		fmt.Println(err)
	}
}

func (p *program) listNamespaces() error {
	namespaces, err := p.namespaceLister.List(labels.NewSelector())
	if err != nil {
		return err
	}
	for _, namespace := range namespaces {
		fmt.Printf("    %s\n", namespace.Name)
	}
	fmt.Println()
	return nil
}

func (p *program) listDNSPods() error {
	const (
		key   = "k8s-app"
		value = "kube-dns"
	)

	r, err := labels.NewRequirement(key, selection.Equals, []string{value})
	if err != nil {
		return err
	}
	pods, err := p.podLister.List(labels.NewSelector().Add(*r))
	if err != nil {
		return err
	}
	for _, pod := range pods {
		fmt.Printf("    name: %s, namespace: %s\n", pod.Name, pod.Namespace)
	}
	fmt.Println()
	return nil
}

func (p *program) createAppNamespace(ctx context.Context) error {
	newNamespace := &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       namespaceKind,
			APIVersion: coreV1APIVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: appNamespace,
		},
	}
	if _, err := p.namespaceClient.Create(ctx, newNamespace, metav1.CreateOptions{}); err != nil {
		if errors.IsAlreadyExists(err) {
			fmt.Printf("Namespace %s already exists.\n\n", appNamespace)
			return nil
		}
		return err
	}
	fmt.Printf("Created namespace %s.\n\n", appNamespace)
	return nil
}

func (p *program) createHelloWorldPod(ctx context.Context) error {
	if _, err := p.podClient.Create(ctx, helloWorldPod(), metav1.CreateOptions{}); err != nil {
		if errors.IsAlreadyExists(err) {
			fmt.Printf("Pod %s already exists.\n\n", appName)
			return nil
		}
		return err
	}
	fmt.Printf("Created pod %s.\n", appName)

	w, err := p.podClient.Watch(ctx, metav1.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(metav1.ObjectNameField, appName).String(),
	})
	if err != nil {
		return err
	}

	observePodCreation(w, creationWaitTimeout)
	return nil
}

func observePodCreation(w watch.Interface, timeout time.Duration) {
	timeoutCh := time.After(timeout)
	defer w.Stop()
	for {
		select {
		case event, ok := <-w.ResultChan():
			if !ok {
				return
			}
			pod, ok := event.Object.(*corev1.Pod)
			if !ok {
				continue
			}
			if pod.Status.Phase == corev1.PodRunning {
				fmt.Printf("Pod status: %s\n\n", corev1.PodRunning)
				return
			}
			if pod.Status.Phase == corev1.PodPending {
				fmt.Println("Pod status:", corev1.PodPending, "...")
			}
		case <-timeoutCh:
			fmt.Println("Time limit exceeded waiting for pod creation", timeout)
			return
		}
	}
}

func (p *program) deleteHelloWorldPod(ctx context.Context) error {
	if err := p.podClient.Delete(ctx, appName, metav1.DeleteOptions{}); err != nil {
		if errors.IsNotFound(err) {
			fmt.Printf("Pod %s doesn't exist.\n", appName)
			return nil
		}
		return err
	}
	fmt.Printf("Deleted Pod %s.\n", appName)
	return nil
}

func helloWorldPod() *corev1.Pod {
	return &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       podKind,
			APIVersion: coreV1APIVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      appName,
			Namespace: appNamespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  appName,
					Image: "us-docker.pkg.dev/google-samples/containers/gke/hello-app:1.0",
					Ports: []corev1.ContainerPort{
						{
							Protocol:      corev1.ProtocolTCP,
							ContainerPort: 8080,
						},
					},
				},
			},
		},
	}
}
