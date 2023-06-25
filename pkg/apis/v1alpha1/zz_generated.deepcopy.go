//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2023 The Kubernetes Authors.

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
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Attach) DeepCopyInto(out *Attach) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Attach.
func (in *Attach) DeepCopy() *Attach {
	if in == nil {
		return nil
	}
	out := new(Attach)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Attach) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AttachConfig) DeepCopyInto(out *AttachConfig) {
	*out = *in
	if in.Containers != nil {
		in, out := &in.Containers, &out.Containers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AttachConfig.
func (in *AttachConfig) DeepCopy() *AttachConfig {
	if in == nil {
		return nil
	}
	out := new(AttachConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AttachSpec) DeepCopyInto(out *AttachSpec) {
	*out = *in
	if in.Attaches != nil {
		in, out := &in.Attaches, &out.Attaches
		*out = make([]AttachConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AttachSpec.
func (in *AttachSpec) DeepCopy() *AttachSpec {
	if in == nil {
		return nil
	}
	out := new(AttachSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterAttach) DeepCopyInto(out *ClusterAttach) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterAttach.
func (in *ClusterAttach) DeepCopy() *ClusterAttach {
	if in == nil {
		return nil
	}
	out := new(ClusterAttach)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterAttach) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterAttachSpec) DeepCopyInto(out *ClusterAttachSpec) {
	*out = *in
	if in.Selector != nil {
		in, out := &in.Selector, &out.Selector
		*out = new(ObjectSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.Attaches != nil {
		in, out := &in.Attaches, &out.Attaches
		*out = make([]AttachConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterAttachSpec.
func (in *ClusterAttachSpec) DeepCopy() *ClusterAttachSpec {
	if in == nil {
		return nil
	}
	out := new(ClusterAttachSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterExec) DeepCopyInto(out *ClusterExec) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterExec.
func (in *ClusterExec) DeepCopy() *ClusterExec {
	if in == nil {
		return nil
	}
	out := new(ClusterExec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterExec) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterExecSpec) DeepCopyInto(out *ClusterExecSpec) {
	*out = *in
	if in.Selector != nil {
		in, out := &in.Selector, &out.Selector
		*out = new(ObjectSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.Execs != nil {
		in, out := &in.Execs, &out.Execs
		*out = make([]ExecTarget, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterExecSpec.
func (in *ClusterExecSpec) DeepCopy() *ClusterExecSpec {
	if in == nil {
		return nil
	}
	out := new(ClusterExecSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterLogs) DeepCopyInto(out *ClusterLogs) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterLogs.
func (in *ClusterLogs) DeepCopy() *ClusterLogs {
	if in == nil {
		return nil
	}
	out := new(ClusterLogs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterLogs) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterLogsSpec) DeepCopyInto(out *ClusterLogsSpec) {
	*out = *in
	if in.Selector != nil {
		in, out := &in.Selector, &out.Selector
		*out = new(ObjectSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.Logs != nil {
		in, out := &in.Logs, &out.Logs
		*out = make([]Log, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterLogsSpec.
func (in *ClusterLogsSpec) DeepCopy() *ClusterLogsSpec {
	if in == nil {
		return nil
	}
	out := new(ClusterLogsSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterMetric) DeepCopyInto(out *ClusterMetric) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterMetric.
func (in *ClusterMetric) DeepCopy() *ClusterMetric {
	if in == nil {
		return nil
	}
	out := new(ClusterMetric)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterMetric) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterMetricSpec) DeepCopyInto(out *ClusterMetricSpec) {
	*out = *in
	if in.Selector != nil {
		in, out := &in.Selector, &out.Selector
		*out = new(ObjectSelectorNames)
		(*in).DeepCopyInto(*out)
	}
	if in.Metrics != nil {
		in, out := &in.Metrics, &out.Metrics
		*out = make([]MetricConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterMetricSpec.
func (in *ClusterMetricSpec) DeepCopy() *ClusterMetricSpec {
	if in == nil {
		return nil
	}
	out := new(ClusterMetricSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterPortForward) DeepCopyInto(out *ClusterPortForward) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterPortForward.
func (in *ClusterPortForward) DeepCopy() *ClusterPortForward {
	if in == nil {
		return nil
	}
	out := new(ClusterPortForward)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterPortForward) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterPortForwardSpec) DeepCopyInto(out *ClusterPortForwardSpec) {
	*out = *in
	if in.Selector != nil {
		in, out := &in.Selector, &out.Selector
		*out = new(ObjectSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.Forwards != nil {
		in, out := &in.Forwards, &out.Forwards
		*out = make([]Forward, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterPortForwardSpec.
func (in *ClusterPortForwardSpec) DeepCopy() *ClusterPortForwardSpec {
	if in == nil {
		return nil
	}
	out := new(ClusterPortForwardSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EnvVar) DeepCopyInto(out *EnvVar) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EnvVar.
func (in *EnvVar) DeepCopy() *EnvVar {
	if in == nil {
		return nil
	}
	out := new(EnvVar)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Exec) DeepCopyInto(out *Exec) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Exec.
func (in *Exec) DeepCopy() *Exec {
	if in == nil {
		return nil
	}
	out := new(Exec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Exec) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExecSpec) DeepCopyInto(out *ExecSpec) {
	*out = *in
	if in.Execs != nil {
		in, out := &in.Execs, &out.Execs
		*out = make([]ExecTarget, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExecSpec.
func (in *ExecSpec) DeepCopy() *ExecSpec {
	if in == nil {
		return nil
	}
	out := new(ExecSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExecTarget) DeepCopyInto(out *ExecTarget) {
	*out = *in
	if in.Containers != nil {
		in, out := &in.Containers, &out.Containers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Local != nil {
		in, out := &in.Local, &out.Local
		*out = new(ExecTargetLocal)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExecTarget.
func (in *ExecTarget) DeepCopy() *ExecTarget {
	if in == nil {
		return nil
	}
	out := new(ExecTarget)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExecTargetLocal) DeepCopyInto(out *ExecTargetLocal) {
	*out = *in
	if in.Envs != nil {
		in, out := &in.Envs, &out.Envs
		*out = make([]EnvVar, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExecTargetLocal.
func (in *ExecTargetLocal) DeepCopy() *ExecTargetLocal {
	if in == nil {
		return nil
	}
	out := new(ExecTargetLocal)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExpressionFromSource) DeepCopyInto(out *ExpressionFromSource) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExpressionFromSource.
func (in *ExpressionFromSource) DeepCopy() *ExpressionFromSource {
	if in == nil {
		return nil
	}
	out := new(ExpressionFromSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FinalizerItem) DeepCopyInto(out *FinalizerItem) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FinalizerItem.
func (in *FinalizerItem) DeepCopy() *FinalizerItem {
	if in == nil {
		return nil
	}
	out := new(FinalizerItem)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Forward) DeepCopyInto(out *Forward) {
	*out = *in
	if in.Ports != nil {
		in, out := &in.Ports, &out.Ports
		*out = make([]int32, len(*in))
		copy(*out, *in)
	}
	if in.Target != nil {
		in, out := &in.Target, &out.Target
		*out = new(ForwardTarget)
		**out = **in
	}
	if in.Command != nil {
		in, out := &in.Command, &out.Command
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Forward.
func (in *Forward) DeepCopy() *Forward {
	if in == nil {
		return nil
	}
	out := new(Forward)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ForwardTarget) DeepCopyInto(out *ForwardTarget) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ForwardTarget.
func (in *ForwardTarget) DeepCopy() *ForwardTarget {
	if in == nil {
		return nil
	}
	out := new(ForwardTarget)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Log) DeepCopyInto(out *Log) {
	*out = *in
	if in.Containers != nil {
		in, out := &in.Containers, &out.Containers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Log.
func (in *Log) DeepCopy() *Log {
	if in == nil {
		return nil
	}
	out := new(Log)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Logs) DeepCopyInto(out *Logs) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Logs.
func (in *Logs) DeepCopy() *Logs {
	if in == nil {
		return nil
	}
	out := new(Logs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Logs) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LogsSpec) DeepCopyInto(out *LogsSpec) {
	*out = *in
	if in.Logs != nil {
		in, out := &in.Logs, &out.Logs
		*out = make([]Log, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LogsSpec.
func (in *LogsSpec) DeepCopy() *LogsSpec {
	if in == nil {
		return nil
	}
	out := new(LogsSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Metric) DeepCopyInto(out *Metric) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Metric.
func (in *Metric) DeepCopy() *Metric {
	if in == nil {
		return nil
	}
	out := new(Metric)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Metric) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetricBucket) DeepCopyInto(out *MetricBucket) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetricBucket.
func (in *MetricBucket) DeepCopy() *MetricBucket {
	if in == nil {
		return nil
	}
	out := new(MetricBucket)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetricConfig) DeepCopyInto(out *MetricConfig) {
	*out = *in
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make([]MetricLabel, len(*in))
		copy(*out, *in)
	}
	if in.Buckets != nil {
		in, out := &in.Buckets, &out.Buckets
		*out = make([]MetricBucket, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetricConfig.
func (in *MetricConfig) DeepCopy() *MetricConfig {
	if in == nil {
		return nil
	}
	out := new(MetricConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetricLabel) DeepCopyInto(out *MetricLabel) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetricLabel.
func (in *MetricLabel) DeepCopy() *MetricLabel {
	if in == nil {
		return nil
	}
	out := new(MetricLabel)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetricSpec) DeepCopyInto(out *MetricSpec) {
	*out = *in
	if in.Metrics != nil {
		in, out := &in.Metrics, &out.Metrics
		*out = make([]MetricConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetricSpec.
func (in *MetricSpec) DeepCopy() *MetricSpec {
	if in == nil {
		return nil
	}
	out := new(MetricSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ObjectSelector) DeepCopyInto(out *ObjectSelector) {
	*out = *in
	if in.MatchNamespaces != nil {
		in, out := &in.MatchNamespaces, &out.MatchNamespaces
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.MatchNames != nil {
		in, out := &in.MatchNames, &out.MatchNames
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ObjectSelector.
func (in *ObjectSelector) DeepCopy() *ObjectSelector {
	if in == nil {
		return nil
	}
	out := new(ObjectSelector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ObjectSelectorNames) DeepCopyInto(out *ObjectSelectorNames) {
	*out = *in
	if in.MatchNames != nil {
		in, out := &in.MatchNames, &out.MatchNames
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ObjectSelectorNames.
func (in *ObjectSelectorNames) DeepCopy() *ObjectSelectorNames {
	if in == nil {
		return nil
	}
	out := new(ObjectSelectorNames)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PortForward) DeepCopyInto(out *PortForward) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PortForward.
func (in *PortForward) DeepCopy() *PortForward {
	if in == nil {
		return nil
	}
	out := new(PortForward)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PortForward) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PortForwardSpec) DeepCopyInto(out *PortForwardSpec) {
	*out = *in
	if in.Forwards != nil {
		in, out := &in.Forwards, &out.Forwards
		*out = make([]Forward, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PortForwardSpec.
func (in *PortForwardSpec) DeepCopy() *PortForwardSpec {
	if in == nil {
		return nil
	}
	out := new(PortForwardSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SelectorRequirement) DeepCopyInto(out *SelectorRequirement) {
	*out = *in
	if in.Values != nil {
		in, out := &in.Values, &out.Values
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SelectorRequirement.
func (in *SelectorRequirement) DeepCopy() *SelectorRequirement {
	if in == nil {
		return nil
	}
	out := new(SelectorRequirement)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Stage) DeepCopyInto(out *Stage) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Stage.
func (in *Stage) DeepCopy() *Stage {
	if in == nil {
		return nil
	}
	out := new(Stage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Stage) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StageDelay) DeepCopyInto(out *StageDelay) {
	*out = *in
	if in.DurationMilliseconds != nil {
		in, out := &in.DurationMilliseconds, &out.DurationMilliseconds
		*out = new(int64)
		**out = **in
	}
	if in.DurationFrom != nil {
		in, out := &in.DurationFrom, &out.DurationFrom
		*out = new(ExpressionFromSource)
		**out = **in
	}
	if in.JitterDurationMilliseconds != nil {
		in, out := &in.JitterDurationMilliseconds, &out.JitterDurationMilliseconds
		*out = new(int64)
		**out = **in
	}
	if in.JitterDurationFrom != nil {
		in, out := &in.JitterDurationFrom, &out.JitterDurationFrom
		*out = new(ExpressionFromSource)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StageDelay.
func (in *StageDelay) DeepCopy() *StageDelay {
	if in == nil {
		return nil
	}
	out := new(StageDelay)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StageEvent) DeepCopyInto(out *StageEvent) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StageEvent.
func (in *StageEvent) DeepCopy() *StageEvent {
	if in == nil {
		return nil
	}
	out := new(StageEvent)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StageFinalizers) DeepCopyInto(out *StageFinalizers) {
	*out = *in
	if in.Add != nil {
		in, out := &in.Add, &out.Add
		*out = make([]FinalizerItem, len(*in))
		copy(*out, *in)
	}
	if in.Remove != nil {
		in, out := &in.Remove, &out.Remove
		*out = make([]FinalizerItem, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StageFinalizers.
func (in *StageFinalizers) DeepCopy() *StageFinalizers {
	if in == nil {
		return nil
	}
	out := new(StageFinalizers)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StageNext) DeepCopyInto(out *StageNext) {
	*out = *in
	if in.Event != nil {
		in, out := &in.Event, &out.Event
		*out = new(StageEvent)
		**out = **in
	}
	if in.Finalizers != nil {
		in, out := &in.Finalizers, &out.Finalizers
		*out = new(StageFinalizers)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StageNext.
func (in *StageNext) DeepCopy() *StageNext {
	if in == nil {
		return nil
	}
	out := new(StageNext)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StageResourceRef) DeepCopyInto(out *StageResourceRef) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StageResourceRef.
func (in *StageResourceRef) DeepCopy() *StageResourceRef {
	if in == nil {
		return nil
	}
	out := new(StageResourceRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StageSelector) DeepCopyInto(out *StageSelector) {
	*out = *in
	if in.MatchLabels != nil {
		in, out := &in.MatchLabels, &out.MatchLabels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.MatchAnnotations != nil {
		in, out := &in.MatchAnnotations, &out.MatchAnnotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.MatchExpressions != nil {
		in, out := &in.MatchExpressions, &out.MatchExpressions
		*out = make([]SelectorRequirement, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StageSelector.
func (in *StageSelector) DeepCopy() *StageSelector {
	if in == nil {
		return nil
	}
	out := new(StageSelector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StageSpec) DeepCopyInto(out *StageSpec) {
	*out = *in
	out.ResourceRef = in.ResourceRef
	if in.Selector != nil {
		in, out := &in.Selector, &out.Selector
		*out = new(StageSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.Delay != nil {
		in, out := &in.Delay, &out.Delay
		*out = new(StageDelay)
		(*in).DeepCopyInto(*out)
	}
	in.Next.DeepCopyInto(&out.Next)
	if in.ImmediateNextStage != nil {
		in, out := &in.ImmediateNextStage, &out.ImmediateNextStage
		*out = new(bool)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StageSpec.
func (in *StageSpec) DeepCopy() *StageSpec {
	if in == nil {
		return nil
	}
	out := new(StageSpec)
	in.DeepCopyInto(out)
	return out
}
