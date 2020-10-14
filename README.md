# ironbird

![An Iron Bird](https://www.aerospacetestinginternational.com/wp-content/uploads/2019/04/Gulfstream_Iron-Bird_2-702x336.jpg)

> An [iron bird](https://en.wikipedia.org/wiki/Iron_bird_(aviation)) is a ground-based test device used for prototyping and integrating aircraft systems during the development of new aircraft designs.

Integration-tests Concourse tasks using `fly execute`, using a YAML test definition format.

## Usage

```terminal
$ ginkgo -- --specs some_spec.yml --target eb [--timeout-factor <int>]
```

Where `timeout-factor` multiplies the default or provided timeouts for execution. Useful if your Concourse is slower than the person who wrote the spec.

Or alternatively

```terminal
$ ginkgo build .
$ ./ironbird-test --specs some_spec.yml  --target eb
```

## Spec Format

See example `*_spec.yml` files because I'm still figuring it out. In the mean time, here's a copypasta that's bound to be out of date:

```yaml
---
# Task config file
config: existing_file_write.yml
cases:
# Each 'when' maps to a `fly execute` invocation
- when: modifier is specified
  # timeout for fly execute
  within: 1m30s
  it:
    # Expected exit code of fly execute
    exits: 0
    # Ordered list of things to expect on STDOUT
    says: [something printed to STDOUT]
    # Define outputs to pull down
    has_outputs:
      - name: output
        for_which:
          # The following bash will be executed and asserted against
          - { bash: "stat existing", exits: 0, says: "4096 0 0 existing" }
          - { bash: "stat modified", exits: 0 }
    # Inputs needed for this test
    has_inputs:
      - name: input
        # Optionally specify a base directory
        from: fixtures/existing_file
        # Optionally run a script in that dir to set it up
        setup: |
          echo foo > modified
  # Param values provided to the task
  params:
    CONTENTS: mycontents
    FILENAME: myfile
```

