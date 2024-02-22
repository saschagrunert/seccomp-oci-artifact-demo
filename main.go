package main

import (
	. "github.com/saschagrunert/demo"
)

func main() {
	d := New()
	d.Add(demo(), "run", "Run the demo")
	d.Run()
}

func demo() *Run {
	r := NewRun(
		"CRI-O - seccomp profiles as OCI artifacts",
		"This demo outlines how to utilize seccomp",
		"profiles as OCI artifacts in the upcoming",
		"CRI-O v1.30.0 release.",
	)

	r.Step(S(
		"CRI-O has to be configured to allow the",
		"annotation `seccomp-profile.kubernetes.cri-o.io`",
		"for the used runtime",
	), S(
		"cat 10-runtimes.conf",
	))

	r.Step(S(
		"We disable the `seccomp_use_default_when_empty` feature",
		"to run containers as `Unconfined` per default. This means",
		"that a workload like this",
	), S(
		"cat pod-unconfined.yml",
	))

	r.Step(nil, S("kubectl apply -f pod-unconfined.yml"))

	r.Step(S(
		"Will run without any apply seccomp profile",
	), S(
		"export CONTAINER_ID=$(sudo crictl ps --name container -q) &&",
		"sudo crictl inspect $CONTAINER_ID | jq .info.runtimeSpec.linux.seccomp",
	))

	r.Step(nil, S("kubectl delete -f pod-unconfined.yml"))

	r.Step(S(
		"We can now use the annotation `seccomp-profile.kubernetes.cri-o.io/POD`",
		"to apply a seccomp profile to a whole pod.",
	), S(
		"cat pod.yml",
	))

	r.Step(nil, S("kubectl apply -f pod.yml"))

	r.Step(S(
		"The pod should now run using the profile, rather",
		"than running `Unconfined`",
	), S(
		"export CONTAINER_ID=$(sudo crictl ps --name container -q) &&",
		"sudo crictl inspect $CONTAINER_ID | jq .info.runtimeSpec.linux.seccomp | head",
	))

	r.Step(nil, S("kubectl delete -f pod.yml"))

	r.Step(S(
		"It's also possible to apply the profile to a specific container",
		"by using the `/<CONTAINER_NAME>` suffix at the anntation level",
	), S(
		"cat pod-ctr.yml",
	))

	r.Step(nil, S("kubectl apply -f pod-ctr.yml"))

	r.Step(S(
		"Again, the container should now run using the profile, rather",
		"than running `Unconfined`",
	), S(
		"export CONTAINER_ID=$(sudo crictl ps --name container -q) &&",
		"sudo crictl inspect $CONTAINER_ID | jq .info.runtimeSpec.linux.seccomp | head",
	))

	r.Step(nil, S("kubectl delete -f pod-ctr.yml"))

	r.Step(S(
		"It's also possible to link container images to seccomp profiles.",
		"For that, the image has to be built using the annotaion",
		"seccomp-profile.kubernetes.cri-o.io=<IMAGE>",
		"",
		"For example, the image quay.io/crio/nginx-seccomp:generic",
		"has been created that way",
	), S(
		"skopeo inspect --raw docker://quay.io/crio/nginx-seccomp:generic |",
		"jq .annotations",
	))

	r.Step(S(
		"If I now use that image in an pod definition",
	), S(
		"cat pod-img.yml",
	))

	r.Step(nil, S("kubectl apply -f pod-img.yml"))

	r.Step(S(
		"Then, the pod should now run using the profile, too",
	), S(
		"export CONTAINER_ID=$(sudo crictl ps --name container -q) &&",
		"sudo crictl inspect $CONTAINER_ID | jq .info.runtimeSpec.linux.seccomp | head",
	))

	r.Step(nil, S("kubectl delete -f pod.yml"))

	r.Step(S(
		"Image annotations also work using the `/POD` or `/<CONTAINER>` suffix.",
		"",
		"The seccomp profile itself consists of a single layer",
		"containing the seccomp.json file in a tar archive.",
	), S(
		"skopeo inspect --raw docker://quay.io/crio/seccomp:v1 | jq .",
	))

	r.Step(S(
		"This allows managing the whole seccomp profile in containers/storage",
		"as well as stacking profiles together (future work).",
		"",
		"Pushing profiles can be done using ORAS: https://oras.land",
	), nil)

	return r
}
