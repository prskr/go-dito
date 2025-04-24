load("@aspect_bazel_lib//lib:stamping.bzl", "STAMP_ATTRS", "maybe_stamp")

# buildifier: disable=module-docstring
RemoteTag = provider(doc = "", fields = ["value"])

def _impl(ctx):
    stamp = maybe_stamp(ctx)
    print(ctx.attr.value)
    if stamp:
        print(stamp.stable_status_file.path)

image_remote_tag = rule(
    implementation = _impl,
    attrs = dict({
    }, **STAMP_ATTRS),
)
