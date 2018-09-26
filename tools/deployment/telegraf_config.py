
class TelegrafConfigFormatter(object):
    NEWLINE='\n'
    EQUALS='='
    QUOT='"'
    COMMA=','
    ARRAYOPEN='['
    ARRAYCLOSE=']'
    INDENT='  '
    def __init__(self, src = None):
        if src is not None:
            self.result = src.result[:] # [:] to copy
            self._current_indent = src._current_indent
        else:
            self.result = []
            self._current_indent = 0

    def to_string(self):
        return ''.join(self.result)

    def append_section_name(self, name, collection=False, inner=False):
        if not inner:
            self._current_indent = 0
        else:
            for _ in range(0, self._current_indent):
                self.result.append(self.INDENT)

        self._current_indent = self._current_indent + 1 if inner else 1

        self.result.append(self.ARRAYOPEN)
        if collection:
            self.result.append(self.ARRAYOPEN)
        self.result.append(name)
        self.result.append(self.ARRAYCLOSE)
        if collection:
            self.result.append(self.ARRAYCLOSE)
        self.result.append(self.NEWLINE)

    def append_key_value(self, key, value):
        for _ in range(0, self._current_indent):
            self.result.append(self.INDENT)

        self.result.append(key)
        self.result.append(self.EQUALS)
        if self._is_iterable(value):
            self._append_collection(value)
        else:
            self._append_value(value)
        self.result.append(self.NEWLINE)

    def _append_collection(self, value):
        self.result.append(self.ARRAYOPEN)
        if len(value) == 0:
            self.result.append(self.ARRAYCLOSE)
            return

        for el in value[:-1]:
            self._append_value(el)
            self.result.append(self.COMMA)

        self._append_value(value[-1])
        self.result.append(self.ARRAYCLOSE)


    def _append_value(self, value):
        if self._is_number(value):
            self.result.append(str(value))
        elif self._is_boolean(value):
            self.result.append(str(value).lower())
        else:
            self.result.append(self.QUOT)
            self.result.append(value)
            self.result.append(self.QUOT)

    def _is_iterable(self, value):
        return not isinstance(value, str) and hasattr(value, '__contains__')

    def _is_number(self, value):
        if self._is_boolean(value):
            return False
        try:
            float(value)
            return True
        except ValueError:
            return False

    def _is_boolean(self, value):
        return isinstance(value, bool)
