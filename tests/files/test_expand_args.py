from nanodoc.files import expand_args


def test_expand_args_empty():
    # Test with empty list
    expanded_files = expand_args([])
    assert expanded_files == []


def test_expand_args_single_file(tmpdir):
    # Create a test file
    test_file = tmpdir.join("test_file.txt")
    test_file.write("test")

    # Call expand_args with a single file path
    expanded_files = expand_args([str(test_file)])

    # Assert that the file path is returned as a single-item list
    assert expanded_files == [str(test_file)]


def test_expand_args_multiple_files(tmpdir):
    # Create test files
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("test1")
    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("test2")

    # Call expand_args with multiple file paths
    expanded_files = expand_args([str(test_file1), str(test_file2)])

    # Assert that both file paths are returned
    assert str(test_file1) in expanded_files
    assert str(test_file2) in expanded_files
    assert len(expanded_files) == 2


def test_expand_args_directory(tmpdir):
    # Create directory structure
    dir_path = tmpdir.mkdir("test_dir")
    test_file_txt = dir_path.join("test_file.txt")
    test_file_txt.write("test")

    # Call expand_args with a directory path
    expanded_files = expand_args([str(dir_path)])

    # Assert that the directory is expanded to include the file
    assert str(test_file_txt) in expanded_files
    assert len(expanded_files) == 1


def test_expand_args_bundle(tmpdir):
    # Create a bundle file
    bundle_file = tmpdir.join("test_bundle.txt")
    test_file = tmpdir.join("test_file.txt")
    test_file.write("test")
    bundle_file.write(str(test_file))

    # Call expand_args with a bundle file path
    expanded_files = expand_args([str(bundle_file)])

    # Assert that the bundle is expanded to include the file
    assert str(test_file) in expanded_files
    assert len(expanded_files) == 1


def test_expand_args_mixed(tmpdir):
    # Create a mix of files, directories, and bundles
    test_file = tmpdir.join("test_file.txt")
    test_file.write("test")

    dir_path = tmpdir.mkdir("test_dir")
    test_file_in_dir = dir_path.join("test_file_in_dir.txt")
    test_file_in_dir.write("test")

    bundle_file = tmpdir.join("test_bundle.txt")
    test_file_for_bundle = tmpdir.join("test_file_for_bundle.txt")
    test_file_for_bundle.write("test")
    bundle_file.write(str(test_file_for_bundle))

    # Call expand_args with a mix of paths
    expanded_files = expand_args([str(test_file), str(dir_path), str(bundle_file)])

    # Assert that all paths are expanded correctly
    assert str(test_file) in expanded_files
    assert str(test_file_in_dir) in expanded_files
    assert str(test_file_for_bundle) in expanded_files
    assert len(expanded_files) == 3
