<h1>Dependabot Wizard</h1>
<h4>The Dependabot Wizard streamlines the setup of basic Dependabot configurations for your projects.</h4>

# Usage

1. Install the dependabot wizard cli:

```bash
go install github.com/kaypee90/dependabot-wizard/cmd/depbot@latest
```

2. Add go/bin to PATH:

```bash
export PATH="$PATH:/usr/local/go/bin"
```
This assumes that your Go installation is in the default location (/usr/local/go). If it's installed elsewhere, adjust the path accordingly.

3. Verify PATH: You can verify that go/bin is included in your PATH by running:

```bash
echo $PATH
```

4. Navigate to the root directory of your project:

```bash
cd <project>
```

5. Start dependabot wizard:

```bash
depbot
```
![Depbot in action](https://github.com/kaypee90/dependabot-wizard/blob/main/assets/depbot.png)
