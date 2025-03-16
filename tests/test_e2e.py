import os
import subprocess
import sys

from nanodoc.formatting import create_header

# Get the parent directory of the current module
MODULE_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))

# Define sample files relative to the fixtures directory
SAMPLE_FILES = [
    os.path.join(MODULE_DIR, "tests", "fixtures", "cake.txt"),
    os.path.join(MODULE_DIR, "tests", "fixtures", "incident.txt"),
    os.path.join(MODULE_DIR, "tests", "fixtures", "new-telephone.txt"),
]

# Use Python module approach instead of direct script execution
PYTHON_CMD = sys.executable
NANODOC_MODULE = "src.nanodoc"


def test_e2e_with_nn_and_toc():
    """
    End-to-end test: process existing sample files with global line numbering
    and TOC.
    """
    # Verify sample files exist
    for file_path in SAMPLE_FILES:
        assert os.path.isfile(file_path), f"Sample file not found: {file_path}"

    # Run nanodoc with -nn and --toc options on the sample files
    result = subprocess.run(
        [PYTHON_CMD, "-m", NANODOC_MODULE, "-nn", "--toc"] + SAMPLE_FILES,
        capture_output=True,
        text=True,
    )

    # Debug: Print the full output
    print("FULL OUTPUT:")
    print(result.stdout)
    print("END OF OUTPUT")

    assert result.returncode == 0

    # Calculate expected output parts (headers)
    toc_header = create_header("TOC", style="filename")
    cake_header = create_header("cake.txt", style="filename")
    incident_header = create_header("incident.txt", style="filename")
    telephone_header = create_header("new-telephone.txt", style="filename")

    # Count lines in sample files to build expectations
    with open(SAMPLE_FILES[0], "r") as f:
        cake_lines = f.readlines()
    with open(SAMPLE_FILES[1], "r") as f:
        incident_lines = f.readlines()
    with open(SAMPLE_FILES[2], "r") as f:
        telephone_lines = f.readlines()

    # Extract just the important parts of the output for comparison (ignoring logs)
    output_lines = result.stdout.split("\n")
    next((i for i, line in enumerate(output_lines) if toc_header in line), 0)
    actual_output = result.stdout

    # Basic assertions about content structure
    assert toc_header in actual_output, "TOC header not found in output"
    assert cake_header in actual_output, "cake.txt header not found"
    assert incident_header in actual_output, "incident.txt header not found"
    assert telephone_header in actual_output, "new-telephone.txt header not found"

    # Check for line numbering
    lines = actual_output.split("\n")
    for i, line in enumerate(lines):
        if any(f"{n}: " in line for n in range(1, 100)):
            assert True  # Found line numbers
            break
    else:
        assert False, "No line numbers found in output"

    # Check TOC contains expected line numbers
    actual_output.split(cake_header)[0]
    # The TOC section includes the TOC header and the lines between it and the first file header
    # Look for each filename in the TOC section directly
    for i, line in enumerate(output_lines):
        if "cake.txt" in line and "..." in line:
            assert True  # Found cake.txt in TOC
            break
    else:
        assert False, "cake.txt not found in TOC"
    # Check for incident.txt in TOC
    for i, line in enumerate(output_lines):
        if "incident.txt" in line and "..." in line:
            assert True  # Found incident.txt in TOC
            break
    else:
        assert False, "incident.txt not found in TOC"
    # Check for new-telephone.txt in TOC
    for i, line in enumerate(output_lines):
        if "new-telephone.txt" in line and "..." in line:
            assert True  # Found new-telephone.txt in TOC
            break
    else:
        assert False, "new-telephone.txt not found in TOC"

    # Line count sanity check
    numbered_lines = [
        line for line in lines if any(f"{n}: " in line for n in range(1, 100))
    ]
    assert len(numbered_lines) == len(cake_lines) + len(incident_lines) + len(
        telephone_lines
    )


def test_e2e_bundle_with_nn_and_toc(tmpdir):
    """
    End-to-end test: process existing sample files via a bundle file with
    global line numbering and TOC.
    """
    # Create bundle file referencing the sample files
    bundle_file = tmpdir.join("bundle.txt")
    bundle_file.write("\n".join(SAMPLE_FILES))

    # Run nanodoc with bundle file, -nn and --toc options
    result = subprocess.run(
        [PYTHON_CMD, "-m", NANODOC_MODULE, "-nn", "--toc", str(bundle_file)],
        capture_output=True,
        text=True,
    )

    # Debug: Print the full output
    print("BUNDLE FULL OUTPUT:")
    print(result.stdout)
    print("END OF BUNDLE OUTPUT")

    assert result.returncode == 0

    # Calculate expected output parts
    toc_header = create_header("TOC", style="filename")
    cake_header = create_header("cake.txt", style="filename")

    # Extract just the important parts of the output
    output_lines = result.stdout.split("\n")
    next((i for i, line in enumerate(output_lines) if toc_header in line), 0)
    actual_output = result.stdout

    # Basic assertions about content
    assert toc_header in actual_output
    assert cake_header in actual_output
    assert (
        os.path.basename(SAMPLE_FILES[1]) in actual_output
    )  # incident.txt basename in output
    assert (
        os.path.basename(SAMPLE_FILES[2]) in actual_output
    )  # new-telephone.txt basename in output

    # Check for line numbering
    lines = actual_output.split("\n")
    for i, line in enumerate(lines):
        if any(f"{n}: " in line for n in range(1, 100)):
            assert True  # Found line numbers
            break
    else:
        assert False, "No line numbers found in output"

    # Check TOC format has filenames and line numbers
    actual_output.split(cake_header)[0]
    for sample_file in SAMPLE_FILES:
        filename = os.path.basename(sample_file)
        # Look for the filename in the TOC section
        found_entry = False
        for line in output_lines:
            if filename in line and "..." in line and any(c.isdigit() for c in line):
                found_entry = True
                break
        assert found_entry, f"TOC entry format incorrect for {filename}"
