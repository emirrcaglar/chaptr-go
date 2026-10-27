# SPDX-License-Identifier: GPL-3.0-or-later
# Go port tooling, 2026-10-27, by emirrcaglar. See ../NOTICE.

"""Regenerate Go parity fixtures using the sibling Rust implementation.

Run from the module root: python tools/generate_parity.py
Requires Python 3 and rustc, with no third-party dependencies.
"""
import json
import pathlib
import re
import subprocess
import tempfile

GO = pathlib.Path(__file__).resolve().parents[1]
RUST = GO.parent / "rust"


def main():
    inputs = set()
    for path in (RUST / "corpus").glob("*.json"):
        inputs.update(row["input"] for row in json.loads(path.read_text(encoding="utf-8")))
    inputs.update(line.strip() for line in (RUST / "corpus/smoke_novel.txt").read_text(encoding="utf-8").splitlines() if line.strip() and not line.startswith("#"))
    # Include every string literal in the upstream unit tests, including lexer
    # and table tests. These exercise malformed and partial input as well.
    for path in (RUST / "src").rglob("*.rs"):
        tests = path.read_text(encoding="utf-8").split("#[cfg(test)]", 1)
        if len(tests) < 2:
            continue
        for literal in re.findall(r'"(?:[^"\\]|\\.)*"', tests[1]):
            try:
                inputs.add(json.loads(literal))
            except json.JSONDecodeError:
                pass
    inputs.update(["", "[", "(", "{", "Title c4294967296", "Title c42v256", "Title v4294967295", "Title c1.65536", "Title c42v0", "Title {r0}"])
    # Keep compiler output in the workspace, including in restricted environments.
    build = GO / ".build"
    build.mkdir(exist_ok=True)
    with tempfile.TemporaryDirectory(dir=build) as tmp:
        tmp = pathlib.Path(tmp)
        lib = tmp / "libchaptr.rlib"
        subprocess.run(["rustc", "--edition=2024", "--crate-type=lib", "--crate-name=chaptr", str(RUST / "src/lib.rs"), "-o", str(lib)], check=True)
        exe = tmp / "oracle.exe"
        subprocess.run(["rustc", "--edition=2024", str(GO / "tools/oracle.rs"), "--extern", f"chaptr={lib}", "-o", str(exe)], check=True)
        # A NUL separator supports embedded newlines and tabs in test literals.
        result = subprocess.run([str(exe)], input="\0".join(sorted(inputs)), text=True, encoding="utf-8", capture_output=True, check=True)
        rows = [json.loads(line) for line in result.stdout.splitlines()]
        output = GO / "chaptr/testdata/rust_parity.json"
        output.write_text(json.dumps(rows, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
        print(f"Wrote {len(rows)} Rust parity fixtures to {output}")


if __name__ == "__main__":
    main()
