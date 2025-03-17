from nanodoc.core import process_all
from nanodoc.files import get_files_from_args
from nanodoc.formatting import create_header


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
    result = process_all(verified_sources)

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
    result = process_all(verified_sources, line_number_mode="file")

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
    result = process_all(verified_sources, line_number_mode="all")

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
    result = process_all(verified_sources, generate_toc=True)

    # Assert that the TOC is generated and the file content is printed
    assert create_header("TOC", style="filename") in result
    assert "test_file1.txt" in result
    assert "test_file2.txt" in result
    assert "Line 1" in result
    assert "Line 2" in result
    assert "Line 3" in result
    assert "Line 4" in result
