package main

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"

	discoveryv1 "k8s.io/api/discovery/v1"
	discoveryv1beta1 "k8s.io/api/discovery/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog/v2/klogr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

func NewClient() (client.Client, error) {
	scheme := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(scheme)

	ctrl.SetLogger(klogr.New())
	cfg := ctrl.GetConfigOrDie()
	cfg.QPS = 100
	cfg.Burst = 100

	hc, err := rest.HTTPClientFor(cfg)
	if err != nil {
		return nil, err
	}
	mapper, err := apiutil.NewDynamicRESTMapper(cfg, hc)
	if err != nil {
		return nil, err
	}

	return client.New(cfg, client.Options{
		Scheme: scheme,
		Mapper: mapper,
		//Opts: client.WarningHandlerOptions{
		//	SuppressWarnings:   false,
		//	AllowDuplicateLogs: false,
		//},
	})
}

func main() {
	//if err := useGeneratedClient(); err != nil {
	//	panic(err)
	//}
	if err := useKubebuilderClient(); err != nil {
		panic(err)
	}
}

func useGeneratedClient() error {
	fmt.Println("Using Generated client")
	cfg := ctrl.GetConfigOrDie()
	cfg.QPS = 100
	cfg.Burst = 100

	kc, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return err
	}

	var list *discoveryv1.EndpointSliceList
	list, err = kc.DiscoveryV1().EndpointSlices(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, db := range list.Items {
		fmt.Println(client.ObjectKeyFromObject(&db))
	}
	return nil
}

func useKubebuilderClient() error {
	fmt.Println("Using kubebuilder client")
	kc, err := NewClient()
	if err != nil {
		return err
	}

	tm, err := BuildTypeMap(kc, schema.GroupVersionKind{
		Group:   discoveryv1.GroupName,
		Version: "v1",
		Kind:    "EndpointSlice",
	})
	if err != nil {
		return err
	}
	fmt.Println(tm)

	k2, err := kc.RESTMapper().KindFor(schema.GroupVersionResource{
		Group:    discoveryv1.GroupName,
		Version:  "",
		Resource: "endpointslices",
	})
	if err != nil {
		return err
	}
	fmt.Println(k2)

	mappings, err := kc.RESTMapper().RESTMappings(schema.GroupKind{
		Group: discoveryv1.GroupName,
		Kind:  "EndpointSlice",
	}, "v1")
	if err != nil {
		return err
	}
	fmt.Println(len(mappings))

	var list discoveryv1beta1.EndpointSliceList
	err = kc.List(context.TODO(), &list)
	if err != nil {
		return err
	}
	scheme := kc.Scheme()

	mappings2, err := kc.RESTMapper().RESTMappings(schema.GroupKind{
		Group: discoveryv1.GroupName,
		Kind:  "EndpointSlice",
	})
	if err != nil {
		return err
	}
	fmt.Println(len(mappings2))

	for _, db := range list.Items {
		// scheme.ConvertToVersion()

		var epv1 discoveryv1.EndpointSlice
		err := scheme.Convert(&db, &epv1, nil)
		if err != nil {
			panic(err)
		}
		fmt.Println(client.ObjectKeyFromObject(&db))
	}

	return nil
}

func BuildTypeMap(kc client.Client, gvks ...schema.GroupVersionKind) (map[schema.GroupVersionKind]schema.GroupVersionKind, error) {
	tm := map[schema.GroupVersionKind]schema.GroupVersionKind{}

	for _, gvk := range gvks {
		mappings, err := kc.RESTMapper().RESTMappings(gvk.GroupKind())
		if err != nil {
			return nil, err
		}

		var found bool
		for _, mapping := range mappings {
			if mapping.GroupVersionKind == gvk {
				found = true
				break
			}
		}
		if !found {
			for _, mapping := range mappings {

				in, err := kc.Scheme().New(gvk)
				if err != nil {
					return nil, err
				}
				out, err := kc.Scheme().New(mapping.GroupVersionKind)
				if err != nil {
					return nil, err
				}
				if err := kc.Scheme().Convert(in, out, nil); err == nil {
					tm[gvk] = mapping.GroupVersionKind
					found = true
					break
				}
			}
		}
		if !found {
			return nil, fmt.Errorf("type mapping not found for %+v", gvk)
		}
	}

	return tm, nil
}
