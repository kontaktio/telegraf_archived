import string
from kapacitor.udf.agent import Agent, Handler
from kapacitor.udf import udf_pb2

from http_store import HttpStore

import logging
logging.basicConfig(level=logging.DEBUG, format='%(asctime)s %(levelname)s:%(name)s: %(message)s')
logger = logging.getLogger()


class SideloadHandler(Handler):

    def __init__(self, agent, store):
        self._agent = agent
        self._field = None
        self._fieldType = 'string'
        self._tag = None
        self._source = None
        self._sourceVariables = []
        self.store = store

    def info(self):
        response = udf_pb2.Response()
        response.info.wants = udf_pb2.STREAM
        response.info.provides = udf_pb2.STREAM
        response.info.options['field'].valueTypes.append(udf_pb2.STRING)
        response.info.options['fieldType'].valueTypes.append(udf_pb2.STRING)
        response.info.options['tag'].valueTypes.append(udf_pb2.STRING)
        response.info.options['source'].valueTypes.append(udf_pb2.STRING)
        return response

    def init(self, init_req):
        success = True
        msg = ''
        for opt in init_req.options:
            if opt.name == 'field':
                self._field = opt.values[0].stringValue
            elif opt.name == 'fieldType':
                self._fieldType = opt.values[0].stringValue
            elif opt.name == 'tag':
                self._tag = opt.values[0].stringValue
            elif opt.name == 'source':
                self._source = opt.values[0].stringValue
                for (_, name, _, _) in string.Formatter().parse(self._source):
                    self._sourceVariables.append(name)

        if (self._field is None or self._fieldType is None) and self._tag is None:
            success = False
            msg += ' must supply field name or tag name. '
        if self._source is None:
            success = False
            msg += ' must supply source address. '

        response = udf_pb2.Response()
        response.init.success = success
        response.init.error = msg[1:]

        return response

    def snapshot(self):
        response = udf_pb2.Response()
        response.snapshot.snapshot = ''

        return response

    def restore(self, restore_req):
        response = udf_pb2.Response()
        response.restore.success = False
        response.restore.error = 'not implemented'

        return response

    def begin_batch(self, begin_req):
        raise Exception("not supported")

    def point(self, point):
        response = udf_pb2.Response()
        response.point.CopyFrom(point)

        kwargs = {}
        for name in self._sourceVariables:
            if name in point.fieldsString:
                kwargs[name] = point.fieldsString[name]
            elif name in point.fieldsInt:
                kwargs[name] = point.fieldsInt[name]
            elif name in point.fieldsDouble:
                kwargs[name] = point.fieldsDouble[name]

        sideload = self.store.call_get(self._source.format(**kwargs))

        if isinstance(sideload, list):
            if len(sideload) > 0:
                sideload = sideload[0]
            else:
                sideload = {}

        if self._field in sideload:
            field = sideload[self._field]
            if self._fieldType == 'string':
                response.point.fieldsDouble[self._field] = field
            elif self._fieldType == 'int':
                response.point.fieldsInt[self._field] = field
            elif self._fieldType == 'double':
                response.point.fieldsDouble[self._field] = field

        if self._tag in sideload:
            tag = sideload[self._tag]
            response.point.tags[self._tag] = tag

        logger.info('Point processed: %s', response.point)

        self._agent.write_response(response, True)

    def end_batch(self, end_req):
        raise Exception("not supported")


if __name__ == '__main__':
    a = Agent()
    h = SideloadHandler(a, HttpStore())
    a.handler = h

    logger.info("Starting Agent")
    a.start()
    a.wait()
    logger.info("Agent finished")
