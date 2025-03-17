from nanodoc.files import is_bundle_file


def test_is_bundle_file_with_valid_paths(tmpdir):
    # Create a test file that will be referenced in the bundle
    test_file = tmpdir.join("referenced_file.txt")
    test_file.write("Some content")

    # Create a bundle file that references the test file
    bundle_file = tmpdir.join("valid_bundle.txt")
    bundle_file.write(str(test_file) + "\n/another/path.txt")

    # Test that the bundle file is recognized as a bundle
    assert is_bundle_file(str(bundle_file))


def test_is_bundle_file_with_content_file(tmpdir):
    # Create a regular content file
    content_file = tmpdir.join("content_file.txt")
    content_file.write("Line 1\nLine 2\nThis is not a file path")

    # Test that the content file is not recognized as a bundle
    assert not is_bundle_file(str(content_file))


def test_is_bundle_file_with_empty_file(tmpdir):
    # Create an empty file
    empty_file = tmpdir.join("empty_file.txt")
    empty_file.write("")

    # Test that an empty file is not recognized as a bundle
    assert not is_bundle_file(str(empty_file))


def test_is_bundle_file_with_invalid_paths(tmpdir):
    # Create a file with invalid paths
    invalid_bundle = tmpdir.join("invalid_bundle.txt")
    invalid_bundle.write("/nonexistent/file1.txt\n/another/nonexistent/file2.txt")

    # Test that a file with only invalid paths is not recognized as a bundle
    assert not is_bundle_file(str(invalid_bundle))


def test_is_bundle_file_with_mixed_content(tmpdir):
    # Create a test file that will be referenced in the bundle
    test_file = tmpdir.join("referenced_file.txt")
    test_file.write("Some content")

    # Create a file with a mix of valid path and regular content
    mixed_file = tmpdir.join("mixed_file.txt")
    mixed_file.write(str(test_file) + "\nThis is just some text\nNot a file path")

    # Test that a file with mixed content is recognized as a bundle
    # since the first valid line is a file path
    assert is_bundle_file(str(mixed_file))


def test_is_bundle_file_with_nonexistent_file():
    # Test with a non-existent file
    assert not is_bundle_file("/nonexistent/path/to/file.txt")


def test_is_bundle_file_with_directory(tmpdir):
    # Create a directory
    dir_path = tmpdir.mkdir("test_dir")

    # Test that a directory is not recognized as a bundle
    assert not is_bundle_file(str(dir_path))


def test_is_bundle_file_with_comments(tmpdir):
    # Create a bundle file with comments
    bundle_file = tmpdir.join("comments_bundle.txt")
    bundle_file.write("# This is a comment\n" + str(tmpdir.join("referenced_file.txt")))

    # Create the referenced file
    referenced_file = tmpdir.join("referenced_file.txt")
    referenced_file.write("Some content")

    # Test that the bundle file is recognized as a bundle
    assert is_bundle_file(str(bundle_file))


def test_is_bundle_file_with_leading_and_trailing_whitespace(tmpdir):
    # Create a bundle file with leading and trailing whitespace
    bundle_file = tmpdir.join("whitespace_bundle.txt")
    bundle_file.write("   " + str(tmpdir.join("referenced_file.txt")) + "   \n")

    # Create the referenced file
    referenced_file = tmpdir.join("referenced_file.txt")
    referenced_file.write("Some content")

    # Test that the bundle file is recognized as a bundle
    assert is_bundle_file(str(bundle_file))


def test_is_bundle_file_with_multiple_lines_and_invalid_paths(tmpdir):
    # Create a bundle file with multiple lines, including invalid paths
    bundle_file = tmpdir.join("multiple_lines_bundle.txt")
    bundle_file.write(
        "/invalid/path1.txt\n"
        + str(tmpdir.join("referenced_file.txt"))
        + "\n/invalid/path2.txt"
    )

    # Create the referenced file
    referenced_file = tmpdir.join("referenced_file.txt")
    referenced_file.write("Some content")

    # With the new implementation, a file is considered a bundle if any line
    # is a valid file path, so this should return True
    assert is_bundle_file(str(bundle_file))
