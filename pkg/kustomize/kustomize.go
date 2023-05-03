package kustomize

import (
	"fmt"

	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

func Kustomize(resources, kustomization string) (string, error) {
	dir := "/"
	memFS := filesys.MakeFsInMemory()
	memFS.WriteFile(dir+"manifest.yaml", []byte(resources))
	memFS.WriteFile(dir+"kustomization.yaml", []byte(kustomization))

	options := &krusty.Options{
		DoLegacyResourceSort: true,
		LoadRestrictions:     types.LoadRestrictionsNone,
		AddManagedbyLabel:    false,
		DoPrune:              false,
		PluginConfig:         types.DisabledPluginConfig(),
	}

	k := krusty.MakeKustomizer(options)
	resmap, err := k.Run(memFS, dir)

	if err != nil {
		return "", fmt.Errorf("kustomize build failed: %w", err)
	}

	yaml, err := resmap.AsYaml()
	if err != nil {
		return "", fmt.Errorf("kustomize build failed: %w", err)
	}

	return string(yaml), nil
}
