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
|   | --ains      | string | start cell\* of analog input ids | |
|   | --aouts     | string | start cell\* of analog output ids | |
|   | --constants | string | start cell\* of constant definitions | |
|   | --dins      | string | start cell\* of digital input ids | |
|   | --douts     | string | start cell\* of digital output ids | |
|   | --flags     | string | start cell\* of flag ids | |
|   | --gins      | string | start cell\* of group input ids | |
|   | --gouts     | string | start cell\* of group output ids | |
| -h| --help      |        | help for fexcel | |
|   | --noupdate  |        | don't check for fexcel updates | |
|   | --numregs   | string | start cell\* of numeric register ids | |
|   | --offset    | int    | column offset between ids and comments | 1 |
|   | --posregs   | string | start cell\* of position register ids | |
|   | --rins      | string | start cell\* of robot input ids | |
|   | --routs     | string | start cell\* of robot output ids | |
|   | --save      |        | save flagset to config file | |
|   | --sheet     | string | default sheet to look at when unspecified in the start cell\* | "Sheet1" |
|   | --sregs     | string | start cell\* of string register ids | |
|   | --timeout   | int    | timeout value in seconds (default 5) |
|   | --ualms     | string | start cell\* of user alarm ids | |

\**start cell flags can be optionally prefixed with a sheet name that
overrides the default `-sheet` flag. (e.g. `--numregs Data:A2`). They
can also include a custom offset (e.g. `--dins 6:IO:A2) where digital
inputs are located on the "IO" sheet starting at A2 and the comments
are in column G.*

## Details

fexcel assumes that your spreadsheet has indices for a given item that start
in the provided cell, and the comments for that item are in the column you
provided plus the offset flag (default 1).

e.g. in the above usage example, the numeric register ids start in cell A2 with
comments starting in cell B2. Position registers ids start in cell D2 with
comments starting in E2. Digital input ids start in cell A2 on the IO sheet.
