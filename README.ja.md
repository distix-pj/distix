# distix

本リポジトリはRPM系Linuxディストリビューションを対象にSBOMを生成するアプリケーションで，[the distix project](https://github.com/distix-pj)の基幹となるものです．
distix projectは，特にRPM系Linuxディストリビューションを対象に，Linuxディストリビューションやそのパッケージ向けに作成するSBOMやその活用方法などについて考えるプロジェクトで，distixはプロジェクトで生成・活用されるSBOMすべてのジェネレータとして実装されています．
distix projectに関して詳しく知りたい方は，以下のリンクを参照してください．
 - [the distix project](https://github.com/distix-pj)

## Getting Started

0. Setting up your Go env..

1. Clone & move to this repository
```
$ git clone https://github.com/distix-pj/distix.git & cd $_
```

2. Install dependencies

```
$ 
```

3. Run command as script
```
$ go run cmd/distix/main.go --help
```

4. or Build & run binary
```
$ go build 
$ ./distix --help
```

5. or Build using our compose.yml on your docker/podman environment
```
$ podman compose run --rm build 
```

## Type & Concept of SBOM distix generates

distixは，いくつかの種類のSBOMを出力することができ，対応するサブコマンドが実装されています．
その種類とコンセプト，およびサブコマンドの実行例について記載します．

### Package SBOM

RPMパッケージファイル(バイナリ)を対象に実行し，そのパッケージのSBOMを生成することができます．
この際，そのパッケージの依存(requires)および包含(provides)関係が，RPMのCapabilityレベルで記載されます．

Package SBOMは，以下のように`package`サブコマンドを実行することで生成できます．
```
$ distix package --input-file [/path/to/package]
```

### Onesystem SBOM

実際に動作しているRPM系Linuxディストリビューションのrpmdbを対象に実行し，そのシステムを対象にしたSBOMを生成することができます．
コンポーネントはパッケージレベルで記載され，それぞれのパッケージの依存関係も記されています．
現在実装されている，Linuxディストリビューションを対象にSBOMを生成できるソフトウェアとほぼ同等の機能で，本機能はその再実装です．
(依存関係の検出方法が異なる可能性はあります)

Onesystem SBOMは，以下のように`onesystem`サブコマンドを実行することで生成できます．
```
$ distix onesystem 
```

### Distsystem SBOM

実際に動作しているRPM系Linuxディストリビューションのrpmdbを対象に実行し，そのシステムを対象にしたSBOMを生成することができます．
SBOMは，そのシステムを対象にしたSBOMファイル一つ，およびシステムにインストールされているパッケージのSBOMがパッケージの数だけ，それぞれ生成されます．
生成されるSBOMはそれぞれ以下のような特徴があります．
 - システムのSBOMは，各パッケージの情報はコンポーネントとして記述されていますが，その依存関係までは記述されていません．
 - システムのSBOM内でコンポーネントとして記載されている各パッケージの情報には，External Referenceとして後述される各パッケージへの情報が記載されています．
 - 各パッケージのSBOMは，Package SBOMで生成されるものと同様に，RPMのCapabilityレベルでの依存・包含関係が記載されます．
本機能は，RPMやディストリビューションの特性を活用するための新たな試みで，その活用例はdistix projectのプロジェクトページや，その配下の別のリポジトリ(ツール)を参照してください．

Distsystem SBOMは，以下のように`distsystem`サブコマンドを実行することで生成できます．また，パッケージのSBOMはデフォルトではsubcompsディレクトリ内に生成されます．
```
$ distix distsystem
```




