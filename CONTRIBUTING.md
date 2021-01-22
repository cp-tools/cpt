# Contributing to cpt

:confetti_ball::tada: Thank you for your interest in contributing to this project! :tada::confetti_ball:

There are some general rules and advices to keep in mind while making Pull Requests.

> The code is more what you call guidelines, than actual rules - *Cpt. Barbossa*

They are not intended to be strictly followed per se, so use your best judgement in making decisions.
Feel free to propose updates and changes to this document!



## Code of Conduct

Be kind, respectful and considerate towards the community. Take time to help others seeking advice. Any form of harassment will not be tolerated, and will be reported. Read the entire rules of conduct [here](CODE_OF_CONDUCT.md).



## Filling a bug report / feature request

1. Before creating a new issue, please check the existing issues to see if any similar one was already opened. Comment on existing ones, rather than creating duplicate issue reports.
2. If you think you've found a bug, please provide detailed steps of reproduction, the version of this library in use, the browser (and version) being automated, and any other useful metrics.
3. If you'd like to see a feature or enhancement, please open an issue with clear descriptions of what you'll like to have, and how its beneficial to the project.



# Guidelines

1. The framework [spf13/cobra](https://github.com/spf13/cobra) is used to create the command line interface. Refer the corresponding documentation for details on incorporating flags, arguments etc.
2. All configuration management shall be done solely by `pkg/conf` to ensure consistency. If there is something missing, consider adding/porting the functionality to the `conf` package.
3. The subroutine code of each sub command must be present in a nested directory. This is to prevent cyclic dependencies and incidental accessing/modification of global scope variables.
4. Ensure inclusion of (color formatted) verbose messages in the subroutines. The color formatting to follow is (can be tweaked, based on the requirement):
   - **Red** - Fatal error messages. Usually followed by `os.Exit(1)`.
   - **Blue** - General verbose messages are to be printed in this color.
   - **Yellow** - Warning messages, but code execution continues.
   - **Green** - Successful execution of some routine.
5. If a code snippet is used across multiple different sub packages, add it to `utils/utils.go`.
6. **Add suitable comments to your code, to let future reviewers know why a given part of the code is required.**




## ELI5: How do I contribute?

First, you need to fork the repository, prior to submitting PRs. Then clone the fork to your computer:

```bash
git clone https://github.com/your_username/cpt.git
cd cpt
```

It is adviced to create a seperate feature branch (than making changes on the master branch):

```bash
git checkout -b my-feature-branch
```

Stage, commit and push changes using the commands:

```bash
git add .
git commit -m "Description of the changes"
git push origin my-feature-branch
```

Once the code is ready, create the Pull Request on GitHub and mark it for review.

The reviewer(s) might suggest changes that should be done. Once satisfied, the PR will be merged, adding your name to the immortal Contributors Hall of Fame! :confetti_ball::confetti_ball: