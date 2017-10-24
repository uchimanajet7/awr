# awr
Create YAML file of AWS term to be used with textlint library phr.

## Description
Output YAML file which can be used for phr of textlint library.

It is acquiring and using AWS terminology of Japanese in the following Web page.

You can customize the expression by writing it in the setting file.

***see also:***

- アマゾン ウェブ サービス： AWS の用語集 
	- https://docs.aws.amazon.com/ja_jp/general/latest/gr/glos-chap.html

## Features
- It is made by golang so it supports multi form.
- You can control the operation in the setting file.
	- You can customize the format of the string you want to check in phr format.


## Requirement
- Go 1.9+
- Packages in use
	- PuerkitoBio/goquery: A little like that j-thing, only in Go.
		- https://github.com/PuerkitoBio/goquery

## Usage
Just run the only one command.

```	sh
$ ./awr
```

However, setting is necessary to execute.

### Setting Example

1. In the same place as the binary file create execution settings file.

1. Execution settings are done with `config.json` file.

```sh
{
	"URL": "https://docs.aws.amazon.com/ja_jp/general/latest/gr/glos-chap.html",
	"Rules": [
		{
			"Expected": "AWS マネジメントコンソール",
			"Patterns": [
				"AWS マネージメントコンソール",
				"AWS Management Console"
			]
		},
		{
			"Expected": "アマゾン ウェブ サービス",
			"Patterns": [
				"Amazon Web Services",
				"Amazon Web Service"
			]
		}
	]
}
```

- About setting items
	- `URL`: String
		- Specify the URL for acquiring the Japanese version of AWS terminology.
	- `Rules`: Array
		- Specify the character string you want to customize as an array.
		- `Expected`: String
			- Describe the correct character string.
			- It must match exactly what is described in the AWS glossary.
		- `Patterns`: Array
			- Write an array of the character strings you want to match with the check.

#### Please edit the output `aws_words.yml` directly if you want to do fine editing.

## Installation

If you build from source yourself.

```	console
$ go get github.com/uchimanajet7/awr
$ cd $GOPATH/src/github.com/uchimanajet7/awr
$ go build
```

### When using only YAML file
Download `aws_words.yml` from the repository and specify it in the phr configuration file.

- awr/aws_words.yml at master · uchimanajet7/awr 
	- https://github.com/uchimanajet7/awr/blob/master/aws_words.yml

***see also:***

- textlint/textlint: The pluggable natural language linter for text and markdown. 
	- https://github.com/textlint/textlint
 
- azu/textlint-rule-prh: textlint rule for prh.
	- https://github.com/azu/textlint-rule-prh

## Author
[uchimanajet7](https://github.com/uchimanajet7)

- textlintを使ってAWS用語をチェックしてみる #aws #textlint #golang - uchimanajet7のメモ
	- http://uchimanajet7.hatenablog.com/entry/2017/10/24/085901

## Licence
[MIT License](https://github.com/uchimanajet7/awr/blob/master/LICENSE)
