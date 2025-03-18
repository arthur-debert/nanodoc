from nanodoc.v1.core import run_all
from nanodoc.v1.files import get_files_from_args
from nanodoc.v1.formatting import create_header


def test_bundle_integration(tmpdir):
    """Test integration of bundle file processing."""
    # Create test files
    file1 = tmpdir.join("file1.txt")
    file2 = tmpdir.join("file2.txt")
    file1.write("Content 1")
    file2.write("Content 2")

    # Create bundle file
    bundle = tmpdir.join("bundle.txt")
    bundle.write(f"{file1}\n{file2}")

    # Process bundle file
    verified_sources = get_files_from_args([str(bundle)])
    result = run_all(
        verified_sources,
        line_number_mode=None,
        generate_toc=False,
        theme=None,
        show_header=True,
        sequence=None,
        style=None,
        txt_ext=None,
    )

    # Verify output contains both files
    assert "file1.txt" in result
    assert "Content 1" in result
    assert "file2.txt" in result
    assert "Content 2" in result


def test_bundle_with_line_ranges(tmpdir):
    """Test bundle file processing with line ranges."""
    # Create test file
    file1 = tmpdir.join("file1.txt")
    file1.write("Line 1\nLine 2\nLine 3\nLine 4")

    # Create bundle file with line ranges
    bundle = tmpdir.join("bundle.txt")
    bundle.write(f"{file1}:L1-2\n{file1}:L3-4")

    # Process bundle file
    verified_sources = get_files_from_args([str(bundle)])
    result = run_all(
        verified_sources,
        line_number_mode=None,
        generate_toc=False,
        theme=None,
        show_header=True,
        sequence=None,
        style=None,
        txt_ext=None,
    )

    # Verify output contains specified line ranges
    assert "file1.txt" in result
    assert "Line 1" in result
    assert "Line 2" in result
    assert "Line 3" in result
    assert "Line 4" in result


def test_bundle_with_toc(tmpdir):
    """Test bundle file processing with table of contents."""
    # Create test files
    file1 = tmpdir.join("file1.txt")
    file2 = tmpdir.join("file2.txt")
    file1.write("Content 1")
    file2.write("Content 2")

    # Create bundle file
    bundle = tmpdir.join("bundle.txt")
    bundle.write(f"{file1}\n{file2}")

    # Process bundle file with TOC
    verified_sources = get_files_from_args([str(bundle)])
    result = run_all(
        verified_sources,
        line_number_mode=None,
        generate_toc=True,
        theme=None,
        show_header=True,
        sequence=None,
        style=None,
        txt_ext=None,
    )

    # Verify output contains TOC and file contents
    assert "Table of Contents" in result
    assert "file1.txt" in result
    assert "file2.txt" in result
    assert "Content 1" in result
    assert "Content 2" in result


def test_init_bundles_no_line_numbers(tmpdir):
    # Create test files
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("Line 1\nLine 2")
    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("Line 3\nLine 4")
    bundle_file = tmpdir.join("test_bundle.txt")
    bundle_file.write(str(test_file1) + "\n" + str(test_file2))

    # Call init with the bundle file
    # Get verified sources and process them
    verified_sources = get_files_from_args([str(bundle_file)])
    result = run_all(
        verified_sources,
        line_number_mode=None,
        generate_toc=False,
        theme=None,
        show_header=True,
        sequence=None,
        style=None,
        txt_ext=None,
    )

    # Assert that the file content is printed without line numbers
    assert "Line 1" in result
    assert "Line 2" in result
    assert "Line 3" in result
    assert "Line 4" in result
    assert "1:" not in result
    assert "2:" not in result
    assert "3:" not in result
    assert "4:" not in result


def test_init_bundles_file_line_numbers(tmpdir):
    # Create test files
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("Line 1\nLine 2")
    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("Line 3\nLine 4")
    bundle_file = tmpdir.join("test_bundle.txt")
    bundle_file.write(str(test_file1) + "\n" + str(test_file2))

    # Call init with the bundle file and file line numbers
    # Get verified sources and process them with file line numbers
    verified_sources = get_files_from_args([str(bundle_file)])
    result = run_all(
        verified_sources,
        line_number_mode="file",
        generate_toc=False,
        theme=None,
        show_header=True,
        sequence=None,
        style=None,
        txt_ext=None,
    )

    # Assert that the file content is printed with file line numbers
    assert "1: Line 1" in result
    assert "2: Line 2" in result
    assert "1: Line 3" in result
    assert "2: Line 4" in result


def test_init_bundles_all_line_numbers(tmpdir):
    # Create test files
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("Line 1\nLine 2")
    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("Line 3\nLine 4")
    bundle_file = tmpdir.join("test_bundle.txt")
    bundle_file.write(str(test_file1) + "\n" + str(test_file2))

    # Call init with the bundle file and all line numbers
    # Get verified sources and process them with all line numbers
    verified_sources = get_files_from_args([str(bundle_file)])
    result = run_all(
        verified_sources,
        line_number_mode="all",
        generate_toc=False,
        theme=None,
        show_header=True,
        sequence=None,
        style=None,
        txt_ext=None,
    )

    # Assert that the file content is printed with all line numbers
    assert "1: Line 1" in result
    assert "2: Line 2" in result
    assert "3: Line 3" in result
    assert "4: Line 4" in result


def test_init_bundles_toc(tmpdir):
    # Create test files
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("Line 1\nLine 2")
    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("Line 3\nLine 4")
    bundle_file = tmpdir.join("test_bundle.txt")
    bundle_file.write(str(test_file1) + "\n" + str(test_file2))

    # Call init with the bundle file and TOC generation
    # Get verified sources and process them with TOC generation
    verified_sources = get_files_from_args([str(bundle_file)])
    result = run_all(
        verified_sources,
        line_number_mode=None,
        generate_toc=True,
        theme=None,
        show_header=True,
        sequence=None,
        style=None,
        txt_ext=None,
    )

    # Assert that the TOC is generated and the file content is printed
    assert create_header("TOC", style="filename") in result
    assert "test_file1.txt" in result
    assert "test_file2.txt" in result
    assert "Line 1" in result
    assert "Line 2" in result
    assert "Line 3" in result
    assert "Line 4" in result
