# gosshtool
[![Build Status](https://travis-ci.org/scottkiss/gosshtool.svg)](https://travis-ci.org/scottkiss/gosshtool)

gosshtool provide some useful functions for ssh client in golang.implemented using golang.org/x/crypto/ssh.


## supports
* command execution on multiple servers.
* ssh tunnel local port forwarding.
* ssh authentication using private keys or password.


## Installation
```bash
go get github.com/scottkiss/gosshtool
```

## Examples

### command execution

```golang
  import "github.com/scottkiss/gosshtool"

	config := &gosshtool.SSHClientConfig{
		User:     "sam",
		Password: "123456",
		Host:     "serverA", //ip:port
	}
	gosshtool.NewSSHClient(config)

	config2 := &gosshtool.SSHClientConfig{
		User:     "sirk",
		Privatekey: "sshprivatekey",
		Host:     "serverB",
	}
	gosshtool.NewSSHClient(config2)
	stdout, _, err := gosshtool.ExecuteCmd("pwd", "serverA")
	if err != nil {
		t.Error(err)
	}
	t.Log(stdout)

	stdout, _, err = gosshtool.ExecuteCmd("pwd", "serverB")
	if err != nil {
		t.Error(err)
	}
	t.Log(stdout)
  ```

### ssh tunnel port forwarding
```golang

package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/scottkiss/gomagic/dbmagic"
	"github.com/scottkiss/gosshtool"
	//"io/ioutil"
	"log"
)

func dbop() {
	ds := new(dbmagic.DataSource)
	ds.Charset = "utf8"
	ds.Host = "127.0.0.1"
	ds.Port = 9999
	ds.DatabaseName = "test"
	ds.User = "root"
	ds.Password = "password"
	dbm, err := dbmagic.Open("mysql", ds)
	if err != nil {
		log.Fatal(err)
	}
	row := dbm.Db.QueryRow("select name from provinces where id=?", 1)
	var name string
	err = row.Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(name)
	dbm.Close()
}

func main() {
	server := new(gosshtool.LocalForwardServer)
	server.LocalBindAddress = ":9999"
	server.RemoteAddress = "remote.com:3306"
	server.SshServerAddress = "112.224.38.111"
	server.SshUserPassword = "passwd"
	//buf, _ := ioutil.ReadFile("/your/home/path/.ssh/id_rsa")
	//server.SshPrivateKey = string(buf)
	server.SshUserName = "sirk"
	server.Start(dbop)
	defer server.Stop()
}

```

## More Examples
* [sshcmd](https://github.com/scottkiss/sshcmd) simple ssh command line client.
* [gooverssh](https://github.com/scottkiss/gooverssh) port forward server over ssh.

## License
View the [LICENSE](https://github.com/scottkiss/gosshtool/blob/master/LICENSE) file


