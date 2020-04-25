Open and transparent principles are important to ensure transparency and cohesiveness in development. With than in mind, the following are Cirrus' design principles:

## Assume a beginner, not an expert
Deploying a CloudFormation stack from end to end as a beginner is often times a frustrating experience. It's often not evident what went wrong and the raw API is difficult to navigate without heavy prior trial-and-error. Cirrus should wire together APIs in a way that a beginner should easily understand.

## Enforce best practices
CloudFormation doesn't do a good job of enforcing best practices. Cirrus can force best practices on users to ensure users don't fall into common traps. This means at times refusing to honor non-best practice inputs. The following are the current enforced best practices:

* `create-stack`/`update-stack` is not allowed. `create-change-set` and `execute-change-set` are the only ways to execute a stack change.
* `delete-stack` is does not trigger a stack deletion. Instead, it triggers a confirmation before the deletion is to begin
* Users cannot pass ad-hoc parameters or tags (e.g. `cirrus up --parameters MyParameter,MyValue`). Tags/parameters must exist as a json file. These files can be sourced, audited, and verified. Ad-hoc parameters/tags cannot.
* `cfn-lint` should be run against the template. cfn-lint may have some false positives so there's currently an option to disable it.

## Have reasonable defaults
The ideal cirrus experience is only having to provide a stack name for the general use-case. Pathological cases should have escape hatches to satisfy their needs, such as a non-standard template name.

## Be human-readable
CLI tools don't have to be ugly to read. While the output from CloudFormation APIs isn't often mean to be human readable, we can fix that. As an example, users often don't care about a stream of events, they care about the state of their resources at any given point. Cirrus focuses on demystifying complex and tedious streams into simple to consume UI components