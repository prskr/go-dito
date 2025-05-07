load("@aspect_bazel_lib//lib:expand_template.bzl", "expand_template")
load("@gazelle//:def.bzl", "gazelle")
load("@rules_go//go:def.bzl", "go_binary", "go_cross_binary", "go_library", "nogo")
load("@rules_oci//oci:defs.bzl", "oci_image", "oci_image_index", "oci_load", "oci_push")
load("@rules_pkg//pkg:tar.bzl", "pkg_tar")

gazelle(name = "gazelle")

ARCHS = [
    "arm64",
    "amd64",
]

nogo(
    name = "analyzers",
    vet = True,
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_gordonklaus_ineffassign//pkg/ineffassign",
        "@com_github_kisielk_errcheck//errcheck",
        "@com_github_lasiar_canonicalheader//:canonicalheader",
    ],
)

go_binary(
    name = "dito",
    embed = [":go-dito_lib"],
    visibility = ["//visibility:public"],
)

go_library(
    name = "go-dito_lib",
    srcs = ["main.go"],
    importpath = "github.com/prskr/go-dito",
    visibility = ["//visibility:private"],
    deps = ["//cmd"],
)

platform(
    name = "linux_arm64",
    constraint_values = [
        "@platforms//os:linux",
        "@platforms//cpu:arm64",
    ],
)

platform(
    name = "linux_amd64",
    constraint_values = [
        "@platforms//os:linux",
        "@platforms//cpu:x86_64",
    ],
)

[
    go_cross_binary(
        name = "dito_" + arch,
        platform = ":linux_" + arch,
        target = ":dito",
    )
    for arch in ARCHS
]

[
    pkg_tar(
        name = "app_layer_" + arch,
        srcs = [":dito_" + arch],
    )
    for arch in ARCHS
]

[
    oci_image(
        name = "image_" + arch,
        architecture = arch,
        entrypoint = ["/dito"],
        os = "linux",
        tars = [":app_layer_" + arch],
    )
    for arch in ARCHS
]

expand_template(
    name = "stamped",
    out = "_stamped.tags.txt",
    stamp_substitutions = {"dev": "{{STABLE_GIT_TAG}}"},
    template = ["dev"],
)

oci_image_index(
    name = "image_multiarch",
    images = [":image_" + arch for arch in ARCHS],
)

oci_push(
    name = "push",
    image = ":image_multiarch",
    remote_tags = ":stamped",
    repository = "ghcr.io/prskr/go-dito",
)

oci_load(
    name = "load",
    image = ":image_arm64",
    repo_tags = ["go-dito:latest"],
)
