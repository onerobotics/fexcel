# fexcel

      __                  _
     / _|                | |
    | |_ _____  _____ ___| |
    |  _/ _ \ \/ / __/ _ \ |
    | ||  __/>  < (_|  __/ |
    |_| \___/_/\_\___\___|_|


Manage your FANUC robot data with an Excel spreadsheet.

Download the latest release [here](https://github.com/onerobotics/fexcel/releases/latest).

[![Build Status](https://travis-ci.org/onerobotics/fexcel.svg "Travis CI status")](https://travis-ci.org/onerobotics/fexcel)

## Usage

Make sure KAREL is unlocked under `Setup > Host Comm > HTTP`.

    fexcel [flags]
    fexcel [commmand] [flags]

Run `fexcel help` for more information on usage.

## Example

There is a `fexcel compile` example located in the `./example` directory.

## Commands

| Command | Description |
| ------- | ----------- |
| compile | Compile a fexcel source file to a FANUC .ls file |
| create  | Create a spreadsheet based on a target's comments |
| diff    | Compare robot comments to spreadsheet (remote or local) |
| help    | Help about any command |
| set     | Set remote robot comments from spreadsheet    |
| version | Print the version number of fexcel |

## Global Flags

|   | Flag        | Type   | Description | Default |
| - | ----------- | ----   | ----------- | ------- |
|   | --ains      | string | start cell of analog input ids | |
|   | --aouts     | string | start cell of analog output ids | |
|   | --constants | string | start cell of constant definitions | |
|   | --dins      | string | start cell of digital input ids | |
|   | --douts     | string | start cell of digital output ids | |
|   | --flags     | string | start cell of flag ids | |
|   | --gins      | string | start cell of group input ids | |
|   | --gouts     | string | start cell of group output ids | |
| -h| --help      |        | help for fexcel | |
|   | --noupdate  |        | don't check for fexcel updates | |
|   | --numregs   | string | start cell of numeric register ids | |
|   | --offset    | int    | column offset between ids and comments | 1 |
|   | --posregs   | string | start cell of position register ids | |
|   | --rins      | string | start cell of robot input ids | |
|   | --routs     | string | start cell of robot output ids | |
|   | --save      |        | save flagset to config file | |
|   | --sheet     | string | default sheet to look at when unspecified in the start cell | "Sheet1" |
|   | --sregs     | string | start cell of string register ids | |
|   | --timeout   | int    | timeout value in seconds (default 5) |
|   | --ualms     | string | start cell of user alarm ids | |

### Start Cell Flags

Providing an excel cell specification (e.g. `A2`) will tell fexcel that the first
numeric id (except for constant definitions, which expect a string identifier)
will be located in the default sheet at cell `A2`. Comments are in the adjacent
column with the default `--offset` of `1`.

You can optionally prefix a sheet name (e.g. `Data:A2`) to override the default sheet.

A custom offset can also be provided (e.g. `A2{6}` or `Data:A2{6}`, indicating comments
are located in column G.

If your definitions are given in multiple places, you can provide multiple start
cells separated by commas (e.g. `A2,D2,Data:A2`).
