package cli

import (
	"flag"
	"fmt"
	
	"github.com/naotama2002/dalv/internal/version"
)

// CLI はコマンドライン引数を処理するための構造体です
type CLI struct {
	helpFlag    *bool
	versionFlag *bool
	tableFlag   *string
	args        []string
}

// NewCLI は新しいCLIインスタンスを作成します
func NewCLI(args []string) *CLI {
	cli := &CLI{
		args: args,
	}
	
	// フラグの定義
	cli.helpFlag = flag.Bool("help", false, "ヘルプ情報を表示します")
	flag.BoolVar(cli.helpFlag, "h", false, "ヘルプ情報を表示します (短縮形)")
	
	cli.versionFlag = flag.Bool("version", false, "バージョン情報を表示します")
	flag.BoolVar(cli.versionFlag, "v", false, "バージョン情報を表示します (短縮形)")
	
	cli.tableFlag = flag.String("table", "", "作成するテーブル名 (デフォルト: 自動生成)")
	flag.StringVar(cli.tableFlag, "t", "", "作成するテーブル名 (短縮形)")
	
	return cli
}

// Parse はコマンドライン引数を解析します
func (c *CLI) Parse() (s3Path string, tableName string, err error) {
	flag.Parse()
	
	// ヘルプフラグが指定された場合
	if *c.helpFlag {
		c.printHelp()
		return "", "", nil
	}
	
	// バージョンフラグが指定された場合
	if *c.versionFlag {
		c.printVersion()
		return "", "", nil
	}
	
	// S3パスの取得
	args := flag.Args()
	if len(args) < 1 {
		return "", "", fmt.Errorf("S3パスが指定されていません。使用方法: dalv [options] <s3-path>")
	}
	
	// テーブル名の取得
	tableName = *c.tableFlag
	
	return args[0], tableName, nil
}

// printHelp はヘルプ情報を表示します
func (c *CLI) printHelp() {
	fmt.Println("dalv - AWS ALB S3ログをDuckDBでクエリするツール")
	fmt.Println()
	fmt.Println("使用方法: dalv [options] <s3-path>")
	fmt.Println()
	fmt.Println("引数:")
	fmt.Println("  <s3-path>  AWS ALBログが保存されているS3パス")
	fmt.Println("             例: s3://bucket/path/to/logs/*.log.gz")
	fmt.Println()
	fmt.Println("オプション:")
	flag.PrintDefaults()
}

// printVersion はバージョン情報を表示します
func (c *CLI) printVersion() {
	fmt.Println(version.VersionString())
}
