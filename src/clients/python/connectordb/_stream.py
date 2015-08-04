from jsonschema import validate, Draft4Validator
import time
import json
from _cobject import ConnectorObject

class Stream(ConnectorObject):
    def create(self,schema):
        #Given a dict representing the stream's JSON Schema, creates the stream
        Draft4Validator.check_schema(schema)
        self.metadata = self.db.urlpost(self.metaname,schema).json()

    @property
    def nickname(self):
        return self.data["nickname"]

    @nickname.setter
    def nickname(self,value):
        self.set({"nickname": value})

    @property
    def downlink(self):
        v = self.data["downlink"]
        if v is None:
            return False
        return v

    @downlink.setter
    def downlink(self,value):
        self.set({"downlink": value})

    @property
    def ephemeral(self):
        v = self.data["ephemeral"]
        if v is None:
            return False
        return v

    @ephemeral.setter
    def ephemeral(self,value):
        self.set({"ephemeral": value})

    @property
    def schema(self):
        return json.loads(self.data["type"])

    def __len__(self):
        return int(self.db.urlget(self.metaname+"/data?q=length").text)

    def insertMany(self,o,restamp=False):
        #attempt to use websocket if websocket inserts are enabled, but fall back on update if fail
        if self.db.wsinsert:
            if self.db.ws.insert(self.metaname,o):
                return
        if restamp:
            self.db.urlput(self.metaname+"/data",o)
        else:
            self.db.urlpost(self.metaname+"/data",o)

    def insert(self,o):
        self.insertMany([{"d":o}],restamp=True)

    def __getitem__(self,obj):
        #Allows to access the stream's elements as if they were an array
        if isinstance(obj,slice):
            start = obj.start
            if start is None:
                start = 0
            stop = obj.stop
            if stop is None:
                stop = 0
            return self.db.urlget(self.metaname+"/data",{"i1": str(start),"i2": str(stop)}).json()
        else:
            return self.db.urlget(self.metaname+"/data",{"i1":str(obj),"i2":str(obj+1)}).json()[0]

    def __call__(self,t1=None,t2=None,limit=None,transform=None,i1=None,i2=None):
        """call allows to get the data range of the stream by any means, and allowing a transform"""
        params = {}
        if not t1 is None:
            params["t1"] = str(t1)
        if not t2 is None:
            params["t2"] = str(t2)
        if not limit is None:
            params["limit"] = str(limit)
        if not i1 is None or not i2 is None:
            if len(params) > 0:
                raise AssertionError("Can't get stream both by index and by timestamp")

            if not i1 is None:
                params["i1"] = str(i1)
            if not t2 is None:
                params["i2"] = str(i2)
        #ConnectorDB doesn't accept null queries
        if len(params)==0:
            params["i1"]="0"

        if not transform is None:
            params["transform"] = transform

        return self.db.urlget(self.metaname+"/data",params).json()

    def subscribe(self,callback,downlink=False):
        '''Stream subscription is a bit more comples, since a stream can be a downlink and can have substreams
        so we subscribe according to that
        '''

        sname = self.metaname
        if downlink:
            sname += "/downlink"
        self.db.ws.subscribe(sname,callback)

    def unsubscribe(self,downlink=False):
        sname = self.metaname
        if downlink:
            sname += "/downlink"
        self.db.ws.unsubscribe(sname)

    def __repr__(self):
        return "[Stream:%s]"%(self.metaname,)
