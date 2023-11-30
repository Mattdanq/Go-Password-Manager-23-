/*
name: Matthew Danque
UIN: 653044344
UID: mdanq2@uic.edu
Program: HOMEWORK 4
This program is supposed to take commands like:
l -> which shows the full list in a nice way
a -> which adds in an ENTRY for the password manager
r -> which removes an ENTRY for the password manager under certain conditions
x -> which exits the program
during the use of the program, it will also have a secondary txt file named passwordVault.txt
which is the full passwordManager in a txt file
*/
package main

import (
	"fmt"
	"os"
  "errors"
  "bufio"
  "log"
)

type EntrySlice []Entry

type Entry struct {
  user string
  pass string
}

// Global variables are allowed (and encouraged) for this project.
var passMap = make(map[string][]Entry)

//_______________________________________________________________________
// initialize before main()
//_______________________________________________________________________
func init () {
  passMap = make(map[string][]Entry)
}

//_______________________________________________________________________
// find the matching entry slice
//_______________________________________________________________________
func findEntrySlice(site string) (EntrySlice, bool) {
  test := passMap[site]
  if (test != nil) {
    return EntrySlice(test), true
  }
  return nil, false
}

//_______________________________________________________________________
// set the entrySlice for site
//_______________________________________________________________________
func setEntrySlice(site string, entrySlice EntrySlice)  {
  passMap[site] = entrySlice
}

//_______________________________________________________________________
// find
//_______________________________________________________________________
func find(user string, entrySlice EntrySlice) (int, bool) {
  for i := 0; i < len(entrySlice); i++ {
    temp := entrySlice[i]
    if temp.user == user {
      return i, true
    }
  }
  return -1, false
}

//_______________________________________________________________________
// print the list in columns
//_______________________________________________________________________
func pmList() {
  
  siteSize := 0
  userSize := 0
  
  for site, slice := range passMap {
    for i := 0; i < len(slice); i++ {
      
      if (siteSize < len(site)) {
        siteSize = len(site)
      }
      
      if (userSize < len(slice[i].user)) {
        userSize = len(slice[i].user)
      }
      
    }
  }

  fmt.Println("List Format: SITE -> USER -> PASSWORD")
  for site, slice := range passMap {
    for i := 0; i < len(slice); i++ {
      var newSite = site
      for a := siteSize - len(site); a > -1 ; a-- {
        newSite = newSite + " "
      }

      var newUser = slice[i].user
      for a := userSize - len(slice[i].user); a > -1; a-- {
        newUser = newUser + " "
      }
      
      fmt.Println(newSite + " " + newUser + " " + slice[i].pass)
    }
  }
}

//_______________________________________________________________________
//  add an entry if the site, user is not already found
//_______________________________________________________________________
func pmAdd(site, user, password string) {
  //new entry
  newEntry := Entry{user: user}
  newEntry.pass = password
  
  sliceCheck, exists := findEntrySlice(site)
  
  if (exists == false && sliceCheck == nil) { //if new site
    tempSlice := make([]Entry, 1)
    tempSlice[0] = newEntry
    //cleanup for later work
    setEntrySlice(site, tempSlice)
    pmWrite()
    
  } else { //if not new site, new ENTRY

    foundIndex, foundUser := find(user, sliceCheck)
    
    if (foundUser == true && foundIndex > -1) { //check dupe
      err := errors.New("add: duplicate entry")
      fmt.Println(err)
      
    } else { //no dupe and real
      //sets up brand new array/vector thing
      tempSlice := passMap[site]
      tempSlice = append(tempSlice, newEntry)
      //cleanup for later work
      setEntrySlice(site, tempSlice)
      pmWrite()
    }
  }
  
}

//_______________________________________________________________________
// remove by site and user
//_______________________________________________________________________
func pmRemove(site, user string) {
  slice, exists := findEntrySlice(site)
  if (exists == false) { //if site doesn't exist
    err := errors.New("remove: site does not exist")
    fmt.Println(err)
    
  } else { //site exists
    userIndex, userExist := find(user, slice)
    
    if (userExist == false) { //user doesnt exist
      err := errors.New("remove: user does not exist")
      fmt.Println(err)
      
    } else { //we create new slice but replace what is removed
      var newSlice []Entry
      newSlice = slice[:userIndex]
      for i := userIndex + 1; i < len(slice); i++ {
        newSlice = append(newSlice, slice[i])
      }
      passMap[site] = newSlice
    }
  }
}

//_______________________________________________________________________
// remove the whole site if there is a single user at that site
//_______________________________________________________________________
func pmRemoveSite(site string) {
    slice, exists := findEntrySlice(site)
    if (exists == false) { //if site doesn't exist
      err := errors.New("remove: site does not exist")
      fmt.Println(err)
      
    } else { //exists
      if (len(slice) > 1) { //if there are multiple entries(users)
        err := errors.New("remove: tried to remove several users")
        fmt.Println(err)
      } else {
        delete(passMap, site)
        pmWrite()
      }
    }
}


//_______________________________________________________________________
// read the passwordVault
//_______________________________________________________________________
func pmRead() {
  
  fileMake, err := os.Create("passwordVault.txt");
  
  //checks if there is an error at runtime
  if (err != nil) {
    panic(err)
    
  }

  fileMake.WriteString("In the format of SITE -> USER -> PASSWORD\n")
  
}

//_______________________________________________________________________
// write the passwordVault
//_______________________________________________________________________
func pmWrite() {
  fileMake, err := os.Create("passwordVault.txt");

  //checks if there is an error at runtime
  if (err != nil) {
    panic(err)

  }


  siteSize := 0
  userSize := 0
  
  fileMake.WriteString("In the format of SITE -> USER -> PASSWORD\n")
  
  for site, slice := range passMap {
    for i := 0; i < len(slice); i++ {

      if (siteSize < len(site)) {
        siteSize = len(site)
      }

      if (userSize < len(slice[i].user)) {
        userSize = len(slice[i].user)
      }

    }
  }

  for site, slice := range passMap {
    for i := 0; i < len(slice); i++ {
      var newSite = site
      for a := siteSize - len(site); a > -1 ; a-- {
        newSite = newSite + " "
      }

      var newUser = slice[i].user
      for a := userSize - len(slice[i].user); a > -1; a-- {
        newUser = newUser + " "
      }

      fileMake.WriteString(newSite + " " + newUser + " " + slice[i].pass + "\n")
    }
  }
}

//_______________________________________________________________________
// do forever loop reading the following commands
//    l
//    a s u p
//    r s
//    r s u
//    x
//  where l,a,r,x are list, add, remove, and exit
//  and s,u,p are site, user, and password
//_______________________________________________________________________
func loop() {
  pmRead()
  for {
    fmt.Println("Input command:")
    
    var command = ""
    var site = ""
    var user = ""
    var password = ""
    
    //takes initial input
    readerIn := bufio.NewReader(os.Stdin)
    stringIn, err := readerIn.ReadString('\n')
    
    //error on input in
    if err != nil {
        log.Fatal(err)
    }
    //makes sure input is always only 4 parameters and not overflowing
    fmt.Sscan(stringIn, &command, &site, &user, &password);
    fmt.Println(command + " " + site + " " + user + " " + password) 
      
    if command == "l" {
      if (site != "" || user != "" || password != "") {
        err := errors.New("list: too much input")
        fmt.Println(err)
        continue;
      }
      pmList()
      
    } else if command == "a" {
      
      //error checking
      if (site == "" || user == "" || password == "") {
        err := errors.New("add: needs more parameters(site, user, password)")
        fmt.Println(err)
        continue;
      }
      
      pmAdd(site, user, password)
      pmWrite()
    } else if command == "r" {
      
      //error checking
      if (site == "") {
        err := errors.New("remove: needs site")
        fmt.Println(err)
        continue
      } else if (password != "") {
        err := errors.New("remove: unneeded password parameter")
        fmt.Println(err)
        continue
      }

      //checks if overriden method or not
      if (user == "") {
        pmRemoveSite(site)
      } else {
        pmRemove(site, user)
      }
      pmWrite()
    } else if command == "x" {
      
      if (site != "" || user != "" || password != "") {
        err := errors.New("add: needs more parameters(site, user, password)")
        fmt.Println(err)
        continue;
      }
      pmWrite()
      break
    }
  }
}

//_______________________________________________________________________
//  let her rip
//_______________________________________________________________________
func main() {
  loop()
}
