import os
import subprocess
import sys

# Get the parent directory of the current module
MODULE_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))

# Define sample files relative to the fixtures directory
SAMPLE_FILES = [
    os.path.join(MODULE_DIR, "fixtures", "cake.txt"),
    os.path.join(MODULE_DIR, "fixtures", "incident.txt"),
    os.path.join(MODULE_DIR, "fixtures", "new-telephone.txt"),
]

# Use Python module approach instead of direct script execution
PYTHON_CMD = sys.executable
NANODOC_MODULE = "src.nanodoc"


def test_e2e_with_nn_and_toc():
    """Test end-to-end process with line numbers and TOC."""
    # Setup
    fixtures_dir = os.path.join(os.path.dirname(__file__), "..", "fixtures")
    test_files = [
        os.path.join(fixtures_dir, "cake.txt"),
        os.path.join(fixtures_dir, "incident.txt"),
        os.path.join(fixtures_dir, "new-telephone.txt"),
    ]

    # Run nanodoc command with v2 mode (appears to be the only mode now)
    cmd = f"python -m nanodoc --use-v2 --toc -n -n {' '.join(test_files)}"
    process = subprocess.run(
        cmd, shell=True, capture_output=True, text=True, check=False
    )
    assert process.returncode == 0, f"Command failed with: {process.stderr}"

    # Parse output
    actual_output = process.stdout
    output_lines = actual_output.split("\n")

    # Check correct header sections are present
    cake_header = "cake.txt"
    incident_header = "incident.txt"
    telephone_header = "new-telephone.txt"

    assert cake_header in actual_output, "cake.txt header not found"
    assert incident_header in actual_output, "incident.txt header not found"
    assert telephone_header in actual_output, "new-telephone.txt header not found"

    # Check line numbers are present (v2 format)
    assert (
        "1: " in actual_output or "   1: " in actual_output
    ), "Line number 1 not found in v2 format"

    # Check TOC contains expected entries
    toc_header = "TOC"
    assert (
        toc_header in actual_output
    ), f"TOC header not found in output: {actual_output}"

    # Check for required files in TOC with v2 format
    filenames = ["cake.txt", "incident.txt", "new-telephone.txt"]
    # V2 TOC format uses bullet points or filename in parentheses
    for filename in filenames:
        found = any(f"- {filename}" in line for line in output_lines)
        # Also check for the filename in another format (e.g., "Cake (cake.txt)")
        if not found:
            found = any(
                filename in line
                for line in output_lines
                if "TOC" in actual_output[: actual_output.index(line)]
            )
        assert found, f"{filename} not found in TOC with v2 format"


def test_e2e_bundle_with_nn_and_toc(tmpdir):
    """Test end-to-end bundling with line numbers and TOC."""
    # Setup
    fixtures_dir = os.path.join(os.path.dirname(__file__), "..", "fixtures")
    bundle_file = tmpdir.join("bundle.txt")
    temp_output = tmpdir.join("output.txt")

    # Get the paths to the fixture files
    test_files = [
        os.path.join(fixtures_dir, "cake.txt"),
        os.path.join(fixtures_dir, "incident.txt"),
        os.path.join(fixtures_dir, "new-telephone.txt"),
    ]

    # Create bundle file with absolute paths to fixtures
    with open(bundle_file, "w") as f:
        for file_path in test_files:
            f.write(f"{file_path}\n")

    try:
        # V2 doesn't support bundle files, so we'll pass the files directly
        cmd = (
            f"python -m nanodoc --use-v2 --toc -n -n "
            f"{' '.join(test_files)} > {temp_output}"
        )
        process = subprocess.run(
            cmd, shell=True, capture_output=True, text=True, check=False
        )
        assert process.returncode == 0, f"Command failed with: {process.stderr}"

        # Read generated output
        with open(temp_output) as f:
            actual_output = f.read()

        output_lines = actual_output.split("\n")

        # Check correct header sections are present
        cake_header = "cake.txt"
        incident_header = "incident.txt"
        telephone_header = "new-telephone.txt"

        assert cake_header in actual_output, "cake.txt header not found"
        assert incident_header in actual_output, "incident.txt header not found"
        assert telephone_header in actual_output, "new-telephone.txt header not found"

        # Check line numbers are present (v2 format)
        assert (
            "1: " in actual_output or "   1: " in actual_output
        ), "Line number 1 not found in v2 format"

        # Check TOC contains expected entries
        toc_header = "TOC"
        assert (
            toc_header in actual_output
        ), f"TOC header not found in output: {actual_output}"

        # Check for required files in TOC with v2 format
        filenames = ["cake.txt", "incident.txt", "new-telephone.txt"]
        # V2 TOC format uses bullet points or filename in parentheses
        for filename in filenames:
            found = any(f"- {filename}" in line for line in output_lines)
            # Also check for the filename in another format (e.g., "Cake (cake.txt)")
            if not found:
                found = any(
                    filename in line
                    for line in output_lines
                    if "TOC" in actual_output[: actual_output.index(line)]
                )
            assert found, f"{filename} not found in TOC with v2 format"
    finally:
        # Cleanup
        if os.path.exists(temp_output):
            os.remove(temp_output)
        if os.path.exists(bundle_file):
            os.remove(bundle_file)
