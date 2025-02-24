## Saitama kills all processes by name at once, with one punch

<p align="center">
  <img width="240" height="200" src="https://raw.githubusercontent.com/lobocode/saitama/master/img/saitama.png">
</p>

## How to install

```bash
curl -s https://raw.githubusercontent.com/lobocode/saitama/master/saitama-install.sh | sudo bash
```

## How to use

![saitama-terminal.gif](https://raw.githubusercontent.com/lobocode/saitama/master/img/saitama-terminal.gif)

### Commands

- `saitama list` - Lists all processes by name.
- `saitama punch <processname>` - Kills the specified process by name.
- `saitama help` - Displays detailed help information.

### Examples

List all processes by name:
```bash
saitama list
```

Kill a specific process by name:
```bash
saitama punch firefox
```

Display help information:
```bash
saitama help
```

With the integration of Cobra, the commands are now more structured and easier to use. Use `saitama help` to get more detailed information about each command.
