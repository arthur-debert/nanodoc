import os

from nanodoc.files import (
    expand_bundles,
    expand_single_arg,
    get_files_from_args,
    is_bundle_file,
    is_file_path_line,
    is_mixed_content_bundle,
    process_mixed_content_bundle,
    process_traditional_bundle,
)


def test_is_file_path_line(tmpdir):
    # Create a test file
    test_file = tmpdir.join("test_file.txt")
    test_file.write("Test content")

    # Test with valid file path
    assert is_file_path_line(str(test_file)) is True

    # Test with non-existent file
    assert is_file_path_line("/nonexistent/file.txt") is False

    # Test with comment
    assert is_file_path_line("# " + str(test_file)) is False

    # Test with empty line - empty strings should return False
    # since they're not valid file paths
    assert is_file_path_line("") is False
    assert is_file_path_line("   ") is False


def test_is_mixed_content_bundle(tmpdir):
    # Create test files
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("Test content 1")
    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("Test content 2")

    # Test with traditional bundle (only file paths)
    lines = [str(test_file1), str(test_file2)]
    assert is_mixed_content_bundle(lines) is False

    # Test with mixed content bundle (text and file paths)
    lines = ["Some text", str(test_file1), "More text"]
    assert is_mixed_content_bundle(lines) is True

    # Test with comments and empty lines
    lines = ["# Comment", "", str(test_file1), "   ", "More text"]
    assert is_mixed_content_bundle(lines) is True

    # Test with only text (no file paths)
    lines = ["Some text", "More text"]
    assert is_mixed_content_bundle(lines) is False


def test_process_mixed_content_bundle(tmpdir):
    # Create test files
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("File content 1")
    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("File content 2")

    # Test with mixed content
    lines = ["Line 1", str(test_file1), "Line 3", str(test_file2), "Line 5"]
    result = process_mixed_content_bundle(lines)

    # Check that file paths were replaced with their content
    assert "Line 1" in result
    assert "File content 1" in result
    assert "Line 3" in result
    assert "File content 2" in result
    assert "Line 5" in result

    # Check the order is preserved
    lines_result = result.splitlines()
    assert lines_result[0] == "Line 1"
    assert lines_result[1] == "File content 1"
    assert lines_result[2] == "Line 3"
    assert lines_result[3] == "File content 2"
    assert lines_result[4] == "Line 5"


def test_process_traditional_bundle(tmpdir):
    # Create test files
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("Test content 1")
    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("Test content 2")

    # Test with traditional bundle
    lines = ["# Comment", "", str(test_file1), "   ", str(test_file2)]
    result = process_traditional_bundle(lines)

    # Check that only file paths were included
    assert len(result) == 2
    assert str(test_file1) in result
    assert str(test_file2) in result


def test_expand_bundles_traditional(tmpdir):
    # Create a traditional bundle file
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("Test content 1")
    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("Test content 2")

    bundle_file = tmpdir.join("traditional_bundle.txt")
    bundle_file.write(f"{test_file1}\n{test_file2}")

    # Test expanding a traditional bundle
    result = expand_bundles(str(bundle_file))

    # Check that it returns a list of file paths
    assert isinstance(result, list)
    assert len(result) == 2
    assert str(test_file1) in result
    assert str(test_file2) in result


def test_expand_bundles_mixed_content(tmpdir):
    # Create test files
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("File content 1")
    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("File content 2")

    # Create a mixed content bundle file
    bundle_file = tmpdir.join("mixed_bundle.txt")
    bundle_file.write(f"Line 1\n{test_file1}\nLine 3\n{test_file2}\nLine 5")

    # Test expanding a mixed content bundle
    result = expand_bundles(str(bundle_file))

    # Check that it returns a string with file paths replaced by their content
    assert isinstance(result, str)
    assert "Line 1" in result
    assert "File content 1" in result
    assert "Line 3" in result
    assert "File content 2" in result
    assert "Line 5" in result


def test_expand_single_arg_mixed_bundle(tmpdir):
    # Create test files
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("File content 1")

    # Create a mixed content bundle file
    bundle_file = tmpdir.join("mixed_bundle.txt")
    bundle_file.write(f"Line 1\n{test_file1}\nLine 3")

    # Make sure the bundle file is recognized as a bundle file
    assert is_bundle_file(str(bundle_file))

    # Test expanding a mixed content bundle
    result = expand_single_arg(str(bundle_file))

    # Check that it returns a list with a single temporary file
    assert isinstance(result, list)
    assert len(result) == 1

    # Check that the temporary file contains the processed content
    with open(result[0]) as f:
        content = f.read()
        assert "Line 1" in content
        # The file path should be replaced with its content
        assert "File content 1" in content
        # The path should not be in the content
        assert str(test_file1) not in content
        assert "Line 3" in content

    # Clean up the temporary file
    os.unlink(result[0])


def test_integration_get_files_from_args(tmpdir):
    # Create test files
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("File content 1")
    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("File content 2")

    # Create a mixed content bundle file
    bundle_file = tmpdir.join("mixed_bundle.txt")
    bundle_file.write(f"Line 1\n{test_file1}\nLine 3\n{test_file2}\nLine 5")

    # Test getting files from args including a mixed content bundle
    result = get_files_from_args([str(bundle_file)])

    # Check that it returns a list of ContentItems
    assert len(result) == 1

    # The first item should be a ContentItem for the temporary file
    assert result[0].file_path.endswith(".txt")

    # Clean up the temporary file
    os.unlink(result[0].file_path)


def test_real_world_example(tmpdir):
    # Create test files to simulate the example in the user query
    lamb_file = tmpdir.join("lamb.txt")
    lamb_file.write("and the lambs are silent")

    # Create a mixed content bundle file
    bundle_file = tmpdir.join("poem.txt")
    bundle_file.write(
        f"Mary had a little lamb\n{lamb_file}\nHis fleece was white as snow, yeah"
    )

    # Test expanding the bundle
    result = expand_bundles(str(bundle_file))

    # Check that it returns the expected content
    expected = (
        "Mary had a little lamb\n"
        "and the lambs are silent\n"
        "His fleece was white as snow, yeah"
    )
    assert result == expected


def test_inline_file_inclusion(tmpdir):
    # Create test files
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write(
        "This is inline content with\nmultiple lines\nthat should be joined."
    )
    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("More inline content")

    # Create a mixed content bundle file with inline file references
    bundle_file = tmpdir.join("inline_bundle.txt")
    bundle_file.write(
        f"Line 1 with @[{test_file1}] in the middle.\n"
        f"Line 2 with @[{test_file2}] at the end.\n"
        f"Line 3 with multiple references: @[{test_file1}] and @[{test_file2}]."
    )

    # Process the bundle directly with process_mixed_content_bundle
    with open(str(bundle_file)) as f:
        lines = f.read().splitlines()

    result = process_mixed_content_bundle(lines)

    # Check that it returns the expected content with inline file inclusions
    expected_inline1 = (
        "This is inline content with multiple lines that should be joined."
    )
    expected_inline2 = "More inline content"

    assert f"Line 1 with {expected_inline1} in the middle." in result
    assert f"Line 2 with {expected_inline2} at the end." in result
    assert (
        f"Line 3 with multiple references: {expected_inline1} and "
        f"{expected_inline2}."
    ) in result

    # Make sure the original file paths are not in the result
    assert f"@[{test_file1}]" not in result
    assert f"@[{test_file2}]" not in result


def test_real_world_inline_example(tmpdir):
    # Create test files to simulate a real-world example
    lamb_file = tmpdir.join("lamb.txt")
    lamb_file.write("and the lambs are silent")

    quote_file = tmpdir.join("quote.txt")
    quote_file.write("To be or not to be\nThat is the question")

    # Create a mixed content bundle file
    bundle_file = tmpdir.join("mixed_inline.txt")
    bundle_file.write(
        f"Mary had a little lamb\n"
        f"{lamb_file}\n"
        f"His fleece was white as snow, yeah\n"
        f"Shakespeare once wrote: @[{quote_file}]"
    )

    # Test expanding the bundle
    result = expand_bundles(str(bundle_file))

    # Check that it returns the expected content
    expected = (
        "Mary had a little lamb\n"
        "and the lambs are silent\n"
        "His fleece was white as snow, yeah\n"
        "Shakespeare once wrote: To be or not to be That is the question"
    )
    assert result == expected
