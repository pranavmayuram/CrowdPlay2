func changeListener () {
    rdb.Send("SUBSCRIBE", pubSubChannel)
    rdb.Flush()
    for {
        reply, err := rdb.Receive()
        if err != nil {
            return err
        }
        // logic of what to do when receiving a change
    }
}


func modelToJSON () {

}
