import { ImageReference } from "@dacc/oci";
import { State } from "dacc";

async function main() {
    const ref = ImageReference.parse("busybox:uclibc")
    const s = await new State().from(ref.toString())

    s.label({ "build.date": new Date().toISOString() })

    console.log(s.toConfig())
}

void main()