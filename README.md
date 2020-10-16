# ironbird

![An Iron Bird](https://www.aerospacetestinginternational.com/wp-content/uploads/2019/04/Gulfstream_Iron-Bird_2-702x336.jpg)

> An [iron bird](https://en.wikipedia.org/wiki/Iron_bird_(aviation)) is a ground-based test device used for prototyping and integrating aircraft systems during the development of new aircraft designs.

Integration-tests Concourse tasks using `fly execute`, using a YAML test definition format.

## Usage

```terminal
$ ironbird --specs some_spec.yml,some_other_spec.yml --target eb [--timeout-factor <int>]
```

* `--specs` - comma-separated list of spec files (see below)
* `--target` - `fly` target on which to execute the tests, and for which there is already a valid auth token
* `--timeout-factor` multiplies the default or provided timeouts for execution. Useful if your Concourse is slower than that of the person who wrote the spec.

## Why?

Concourse tasks should be tested, and there was no simple, succinct way to test simple tasks. No-one wants to write a hundred lines of Golang to test four lines of Bash.

## Spec Format

See `*_spec.yml` files in `integration` for examples.

```yaml
---
# Task config file (required)
config: existing_file_write.yml
# details of input that task config is normally within (optional)
enclosed_in_input:
  # the name of the input containing the task.yml (optional)
  name: some-repo
  # where the 'root' of the input containing task YAML/script is, relative to this spec file (optional)
  path_relative_to_spec: ../../../
cases:
# Each 'when' maps to a `fly execute` invocation
- when: modifier is specified
  # timeout for fly execute (optional, defaults to 20s)
  within: 1m30s
  it:
    # Expected exit code of fly execute (optional, defaults to 0)
    exits: 0
    # Ordered list of things to expect on STDOUT (optional)
    says: [something printed to STDOUT]
    # Define outputs to pull down (optional)
    has_outputs:
      - name: output
        for_which:
          # The following bash will be executed and asserted against
          - { bash: "stat existing", exits: 0, says: "4096 0 0 existing" }
          - { bash: "stat modified", exits: 0 }
    # Additional inputs needed for this test (optional)
    has_inputs:
      - name: input
        # Dir, relative to this spec file, to use as the basis for the input (optional)
        from: fixtures/existing_file
        # Bash commands to apply to either blank dir, or dir above, before running fly execute (optional)
        setup: |
          echo foo > modified
  # Param values provided to the task (optional)
  params:
    CONTENTS: mycontents
    FILENAME: myfile
```

## Wouldn't it be quicker if the tests ran in Docker?

Yes.
