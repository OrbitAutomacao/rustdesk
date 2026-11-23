"""Prepare the embedded RustDesk client HTML with the Orbit WebSocket shim."""

import argparse
from pathlib import Path
from urllib.request import urlopen


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--origin", default="http://127.0.0.1:21114/webclient/")
    parser.add_argument("--output", required=True, type=Path)
    args = parser.parse_args()

    with urlopen(args.origin, timeout=10) as response:
        source = response.read().decode("utf-8")

    marker = '<script src="/webclient-config/index.js"></script>'
    if source.count(marker) != 1:
        raise ValueError("Unexpected RustDesk web client HTML; shim was not inserted")

    patched = source.replace(
        marker,
        marker + '<script src="/webclient-shim.js?v=2"></script>',
        1,
    )
    args.output.write_text(patched, encoding="utf-8")


if __name__ == "__main__":
    main()
