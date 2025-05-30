# Minecraft Input Plugin

The `minecraft` plugin connects to a Minecraft server using the RCON protocol
to collects scores from the server [scoreboard][].

This plugin is known to support Minecraft Java Edition versions 1.11 - 1.14.
When using an version of Minecraft earlier than 1.13, be aware that the values
for some criterion has changed and may need to be modified.

## Server Setup

Enable [RCON][] on the Minecraft server, add this to your server configuration
in the [server.properties][] file:

```conf
#  Minecraft Plugin

This plugin uses the RCON protocol to collect [statistics](http://minecraft.gamepedia.com/Statistics) from a [scoreboard](http://minecraft.gamepedia.com/Scoreboard) on a
Minecraft server.

To enable [RCON](http://wiki.vg/RCON) on the minecraft server, add this to your server configuration in the `server.properties` file:

```
enable-rcon=true
rcon.password=<your password>
rcon.port=<1-65535>
```

Scoreboard [Objectives][] must be added using the server console for the
plugin to collect.  These can be added in game by players with op status,
from the server console, or over an RCON connection.

When getting started pick an easy to test objective.  This command will add an
objective that counts the number of times a player has jumped:

```sh
/scoreboard objectives add jumps minecraft.custom:minecraft.jump
```

Once a player has triggered the event they will be added to the scoreboard,
you can then list all players with recorded scores:

```sh
/scoreboard players list
```

View the current scores with a command, substituting your player name:

```sh
/scoreboard players list Etho
```

## Global configuration options <!-- @/docs/includes/plugin_config.md -->

In addition to the plugin-specific configuration settings, plugins support
additional global and plugin configuration settings. These settings are used to
modify metrics, tags, and field or create aliases and configure ordering, etc.
See the [CONFIGURATION.md][CONFIGURATION.md] for more details.

[CONFIGURATION.md]: ../../../docs/CONFIGURATION.md

## Configuration

```toml @sample.conf
# Collects scores from a Minecraft server's scoreboard using the RCON protocol
[[inputs.minecraft]]
  ## Address of the Minecraft server.
  # server = "localhost"

  ## Server RCON Port.
  # port = "25575"

  ## Server RCON Password.
  password = ""

  ## Uncomment to remove deprecated metric components.
  # tagdrop = ["server"]
```

## Metrics

- minecraft
  - tags:
    - player
    - port (port of the server)
    - server (hostname:port, deprecated in 1.11; use `source` and `port` tags)
    - source (hostname of the server)
  - fields:
    - `<objective_name>` (integer, count)

## Sample Queries

Get the number of jumps per player in the last hour:

```sql
SELECT SPREAD("jumps") FROM "minecraft" WHERE time > now() - 1h GROUP BY "player"
```

## Example Output

```shell
minecraft,player=notch,source=127.0.0.1,port=25575 jumps=178i 1498261397000000000
minecraft,player=dinnerbone,source=127.0.0.1,port=25575 deaths=1i,jumps=1999i,cow_kills=1i 1498261397000000000
minecraft,player=jeb,source=127.0.0.1,port=25575 d_pickaxe=1i,damage_dealt=80i,d_sword=2i,hunger=20i,health=20i,kills=1i,level=33i,jumps=264i,armor=15i 1498261397000000000
```

[server.properties]: https://minecraft.gamepedia.com/Server.properties
[scoreboard]: http://minecraft.gamepedia.com/Scoreboard
[objectives]: https://minecraft.gamepedia.com/Scoreboard#Objectives
[rcon]: http://wiki.vg/RCON
To create a new scoreboard objective called `jump` on a minecraft server tracking the `stat.jump` criteria, run this command
in the Minecraft console:

`/scoreboard objectives add jump stat.jump`

Stats are collected with the following RCON command, issued by the plugin:

`/scoreboard players list *`

### Configuration:
```
[[inputs.minecraft]]
   # server address for minecraft
   server = "localhost"
   # port for RCON
   port = "25575"
   # password RCON for mincraft server
   password = "replace_me"
```

### Measurements & Fields:

*This plugin uses only one measurement, titled* `minecraft`

- The field name is the scoreboard objective name.
- The field value is the count of the scoreboard objective

- `minecraft`
    - `<objective_name>` (integer, count)

### Tags:

- The `minecraft` measurement:
    - `server`: the Minecraft RCON server
    - `player`: the Minecraft player


### Sample Queries:

Get the number of jumps per player in the last hour:
```
SELECT SPREAD("jump") FROM "minecraft" WHERE time > now() - 1h GROUP BY "player"
```

### Example Output:

```
$ telegraf --input-filter minecraft --test
* Plugin: inputs.minecraft, Collection 1
> minecraft,player=notch,server=127.0.0.1:25575 jumps=178i 1498261397000000000
> minecraft,player=dinnerbone,server=127.0.0.1:25575 deaths=1i,jumps=1999i,cow_kills=1i 1498261397000000000
> minecraft,player=jeb,server=127.0.0.1:25575 d_pickaxe=1i,damage_dealt=80i,d_sword=2i,hunger=20i,health=20i,kills=1i,level=33i,jumps=264i,armor=15i 1498261397000000000
```
