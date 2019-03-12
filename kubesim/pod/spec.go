package pod

import (
	"github.com/cpuguy83/strongerrors"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"

	"github.com/ordovicia/k8s-cluster-simulator/kubesim/util"
)

// spec represents a list of a pod's resource usage spec of each execution phase.
type spec []specPhase

// specPhase represents a pod's resource usage spec of one execution phase.
type specPhase struct {
	seconds       int32
	resourceUsage v1.ResourceList
}

// errInvalidResourceUsageField is returned from parseSpec.
var errInvalidResourceUsageField = errors.New("Invalid spec.resoruceUsage field")

// parseSpec parses the pod's "spec" annotation into spec.
// Returns error if the "spec" annotation does not exist, or the parsing failes.
func parseSpec(pod *v1.Pod) (spec, error) {
	specAnnot, ok := pod.ObjectMeta.Annotations["simSpec"]
	if !ok {
		return nil, strongerrors.InvalidArgument(errors.Errorf("simSpec annotation not defined"))
	}

	return parseSpecYAML(specAnnot)
}

// parseSpec parses the YAML into spec.
// Returns error if the YAML is invalid.
func parseSpecYAML(specYAML string) (spec, error) {
	type specPhaseYAML struct {
		Seconds       int32                      `yaml:"seconds"`
		ResourceUsage map[v1.ResourceName]string `yaml:"resourceUsage"`
	}

	specUnmarshalled := []specPhaseYAML{}
	if err := yaml.Unmarshal([]byte(specYAML), &specUnmarshalled); err != nil {
		return nil, err
	}

	spec := spec{}
	for _, phase := range specUnmarshalled {
		if phase.ResourceUsage == nil {
			return nil, errInvalidResourceUsageField
		}

		resourceUsage, err := util.BuildResourceList(phase.ResourceUsage)
		if err != nil {
			return nil, err
		}
		spec = append(spec, specPhase{
			seconds:       phase.Seconds,
			resourceUsage: resourceUsage,
		})
	}

	return spec, nil
}
