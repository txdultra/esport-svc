```
package main

import (
    "log"
    "github.com/9nmy/SequoiaDB"
)

func main(){
    session,rc := SequoiaDB.Connect("localhost", "11810", "", "");
    if rc != 0 {
        log.Fatal(rc)
    }

    cursor := session.Query("select _id as id,username,password from nmy.user")
    for d,e := cursor.Next(); e != 0 {
        log.Println(d)
    }
    
    session.Disconnect()
}

```
