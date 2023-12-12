//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

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

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	json "encoding/json"

	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Component) DeepCopyInto(out *Component) {
	*out = *in
	if in.Links != nil {
		in, out := &in.Links, &out.Links
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Command != nil {
		in, out := &in.Command, &out.Command
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Args != nil {
		in, out := &in.Args, &out.Args
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Ports != nil {
		in, out := &in.Ports, &out.Ports
		*out = make([]Port, len(*in))
		copy(*out, *in)
	}
	if in.Envs != nil {
		in, out := &in.Envs, &out.Envs
		*out = make([]Env, len(*in))
		copy(*out, *in)
	}
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]Volume, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Metric != nil {
		in, out := &in.Metric, &out.Metric
		*out = new(ComponentMetric)
		**out = **in
	}
	if in.MetricsDiscovery != nil {
		in, out := &in.MetricsDiscovery, &out.MetricsDiscovery
		*out = new(ComponentMetric)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Component.
func (in *Component) DeepCopy() *Component {
	if in == nil {
		return nil
	}
	out := new(Component)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ComponentMetric) DeepCopyInto(out *ComponentMetric) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ComponentMetric.
func (in *ComponentMetric) DeepCopy() *ComponentMetric {
	if in == nil {
		return nil
	}
	out := new(ComponentMetric)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ComponentPatches) DeepCopyInto(out *ComponentPatches) {
	*out = *in
	if in.ExtraArgs != nil {
		in, out := &in.ExtraArgs, &out.ExtraArgs
		*out = make([]ExtraArgs, len(*in))
		copy(*out, *in)
	}
	if in.ExtraVolumes != nil {
		in, out := &in.ExtraVolumes, &out.ExtraVolumes
		*out = make([]Volume, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.ExtraEnvs != nil {
		in, out := &in.ExtraEnvs, &out.ExtraEnvs
		*out = make([]Env, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ComponentPatches.
func (in *ComponentPatches) DeepCopy() *ComponentPatches {
	if in == nil {
		return nil
	}
	out := new(ComponentPatches)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Env) DeepCopyInto(out *Env) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Env.
func (in *Env) DeepCopy() *Env {
	if in == nil {
		return nil
	}
	out := new(Env)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExtraArgs) DeepCopyInto(out *ExtraArgs) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExtraArgs.
func (in *ExtraArgs) DeepCopy() *ExtraArgs {
	if in == nil {
		return nil
	}
	out := new(ExtraArgs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KwokConfiguration) DeepCopyInto(out *KwokConfiguration) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Options.DeepCopyInto(&out.Options)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KwokConfiguration.
func (in *KwokConfiguration) DeepCopy() *KwokConfiguration {
	if in == nil {
		return nil
	}
	out := new(KwokConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KwokConfiguration) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KwokConfigurationOptions) DeepCopyInto(out *KwokConfigurationOptions) {
	*out = *in
	if in.EnableCRDs != nil {
		in, out := &in.EnableCRDs, &out.EnableCRDs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ManageAllNodes != nil {
		in, out := &in.ManageAllNodes, &out.ManageAllNodes
		*out = new(bool)
		**out = **in
	}
	if in.EnableCNI != nil {
		in, out := &in.EnableCNI, &out.EnableCNI
		*out = new(bool)
		**out = **in
	}
	if in.EnableDebuggingHandlers != nil {
		in, out := &in.EnableDebuggingHandlers, &out.EnableDebuggingHandlers
		*out = new(bool)
		**out = **in
	}
	if in.EnableContentionProfiling != nil {
		in, out := &in.EnableContentionProfiling, &out.EnableContentionProfiling
		*out = new(bool)
		**out = **in
	}
	if in.EnableProfilingHandler != nil {
		in, out := &in.EnableProfilingHandler, &out.EnableProfilingHandler
		*out = new(bool)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KwokConfigurationOptions.
func (in *KwokConfigurationOptions) DeepCopy() *KwokConfigurationOptions {
	if in == nil {
		return nil
	}
	out := new(KwokConfigurationOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KwokctlConfiguration) DeepCopyInto(out *KwokctlConfiguration) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Options.DeepCopyInto(&out.Options)
	if in.Components != nil {
		in, out := &in.Components, &out.Components
		*out = make([]Component, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.ComponentsPatches != nil {
		in, out := &in.ComponentsPatches, &out.ComponentsPatches
		*out = make([]ComponentPatches, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KwokctlConfiguration.
func (in *KwokctlConfiguration) DeepCopy() *KwokctlConfiguration {
	if in == nil {
		return nil
	}
	out := new(KwokctlConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KwokctlConfiguration) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KwokctlConfigurationOptions) DeepCopyInto(out *KwokctlConfigurationOptions) {
	*out = *in
	if in.EnableCRDs != nil {
		in, out := &in.EnableCRDs, &out.EnableCRDs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Runtimes != nil {
		in, out := &in.Runtimes, &out.Runtimes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.SecurePort != nil {
		in, out := &in.SecurePort, &out.SecurePort
		*out = new(bool)
		**out = **in
	}
	if in.QuietPull != nil {
		in, out := &in.QuietPull, &out.QuietPull
		*out = new(bool)
		**out = **in
	}
	if in.DisableKubeScheduler != nil {
		in, out := &in.DisableKubeScheduler, &out.DisableKubeScheduler
		*out = new(bool)
		**out = **in
	}
	if in.DisableKubeControllerManager != nil {
		in, out := &in.DisableKubeControllerManager, &out.DisableKubeControllerManager
		*out = new(bool)
		**out = **in
	}
	if in.EnableMetricsServer != nil {
		in, out := &in.EnableMetricsServer, &out.EnableMetricsServer
		*out = new(bool)
		**out = **in
	}
	if in.KubeAuthorization != nil {
		in, out := &in.KubeAuthorization, &out.KubeAuthorization
		*out = new(bool)
		**out = **in
	}
	if in.KubeAdmission != nil {
		in, out := &in.KubeAdmission, &out.KubeAdmission
		*out = new(bool)
		**out = **in
	}
	if in.KubeApiserverCertSANs != nil {
		in, out := &in.KubeApiserverCertSANs, &out.KubeApiserverCertSANs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.DisableQPSLimits != nil {
		in, out := &in.DisableQPSLimits, &out.DisableQPSLimits
		*out = new(bool)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KwokctlConfigurationOptions.
func (in *KwokctlConfigurationOptions) DeepCopy() *KwokctlConfigurationOptions {
	if in == nil {
		return nil
	}
	out := new(KwokctlConfigurationOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KwokctlConfigurationStatus) DeepCopyInto(out *KwokctlConfigurationStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KwokctlConfigurationStatus.
func (in *KwokctlConfigurationStatus) DeepCopy() *KwokctlConfigurationStatus {
	if in == nil {
		return nil
	}
	out := new(KwokctlConfigurationStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KwokctlResource) DeepCopyInto(out *KwokctlResource) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Parameters != nil {
		in, out := &in.Parameters, &out.Parameters
		*out = make(json.RawMessage, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KwokctlResource.
func (in *KwokctlResource) DeepCopy() *KwokctlResource {
	if in == nil {
		return nil
	}
	out := new(KwokctlResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KwokctlResource) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Port) DeepCopyInto(out *Port) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Port.
func (in *Port) DeepCopy() *Port {
	if in == nil {
		return nil
	}
	out := new(Port)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Volume) DeepCopyInto(out *Volume) {
	*out = *in
	if in.ReadOnly != nil {
		in, out := &in.ReadOnly, &out.ReadOnly
		*out = new(bool)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Volume.
func (in *Volume) DeepCopy() *Volume {
	if in == nil {
		return nil
	}
	out := new(Volume)
	in.DeepCopyInto(out)
	return out
}
