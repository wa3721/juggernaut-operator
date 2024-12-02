//go:build !ignore_autogenerated

/*
Copyright 2024.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Juggernaut) DeepCopyInto(out *Juggernaut) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Juggernaut.
func (in *Juggernaut) DeepCopy() *Juggernaut {
	if in == nil {
		return nil
	}
	out := new(Juggernaut)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Juggernaut) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JuggernautConfig) DeepCopyInto(out *JuggernautConfig) {
	*out = *in
	out.Overwrite = in.Overwrite
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JuggernautConfig.
func (in *JuggernautConfig) DeepCopy() *JuggernautConfig {
	if in == nil {
		return nil
	}
	out := new(JuggernautConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JuggernautConfigmap) DeepCopyInto(out *JuggernautConfigmap) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JuggernautConfigmap.
func (in *JuggernautConfigmap) DeepCopy() *JuggernautConfigmap {
	if in == nil {
		return nil
	}
	out := new(JuggernautConfigmap)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JuggernautList) DeepCopyInto(out *JuggernautList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Juggernaut, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JuggernautList.
func (in *JuggernautList) DeepCopy() *JuggernautList {
	if in == nil {
		return nil
	}
	out := new(JuggernautList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *JuggernautList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JuggernautService) DeepCopyInto(out *JuggernautService) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JuggernautService.
func (in *JuggernautService) DeepCopy() *JuggernautService {
	if in == nil {
		return nil
	}
	out := new(JuggernautService)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JuggernautSpec) DeepCopyInto(out *JuggernautSpec) {
	*out = *in
	in.Resources.DeepCopyInto(&out.Resources)
	out.Service = in.Service
	out.Config = in.Config
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JuggernautSpec.
func (in *JuggernautSpec) DeepCopy() *JuggernautSpec {
	if in == nil {
		return nil
	}
	out := new(JuggernautSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JuggernautStatus) DeepCopyInto(out *JuggernautStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JuggernautStatus.
func (in *JuggernautStatus) DeepCopy() *JuggernautStatus {
	if in == nil {
		return nil
	}
	out := new(JuggernautStatus)
	in.DeepCopyInto(out)
	return out
}
