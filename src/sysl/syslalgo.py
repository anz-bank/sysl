# Copyright 2016 The Sysl Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License."""Super smart code writer."""

"""sysl algorithms"""

def enumerate_calls(stmts):  # pylint: disable=too-many-branches
  """Enumerate all calls under stmts, in the form (parent_stmt, call)."""
  for stmt in stmts:
    if stmt.HasField('call'):
      yield (stmt, stmt.call)
    elif stmt.HasField('action'):
      pass
    elif stmt.HasField('cond'):
      for call in enumerate_calls(stmt.cond.stmt):
        yield call
    elif stmt.HasField('loop'):
      for call in enumerate_calls(stmt.loop.stmt):
        yield call
    elif stmt.HasField('loop_n'):
      for call in enumerate_calls(stmt.loop_n.stmt):
        yield call
    elif stmt.HasField('foreach'):
      for call in enumerate_calls(stmt.foreach.stmt):
        yield call
    elif stmt.HasField('group'):
      for call in enumerate_calls(stmt.group.stmt):
        yield call
    elif stmt.HasField('alt'):
      for choice in stmt.alt.choice:
        for call in enumerate_calls(choice.stmt):
          yield call
    elif stmt.HasField('ret'):
      pass
    else:
      raise Exception('No statement!', stmt.WhichOneof('stmt'))


# TODO: Require consistency along different branches.
def return_payload(stmts):
  # pylint: disable=too-many-branches,too-many-return-statements
  """Compute payload returned by stmts."""
  for stmt in stmts:
    if stmt.HasField('call'):
      pass
    elif stmt.HasField('action'):
      pass
    elif stmt.HasField('cond'):
      payload = return_payload(stmt.cond.stmt)
      if payload:
        return payload
    elif stmt.HasField('loop'):
      payload = return_payload(stmt.loop.stmt)
      if payload:
        return payload
    elif stmt.HasField('loop_n'):
      payload = return_payload(stmt.loop_n.stmt)
      if payload:
        return payload
    elif stmt.HasField('foreach'):
      payload = return_payload(stmt.foreach.stmt)
      if payload:
        return payload
    elif stmt.HasField('group'):
      payload = return_payload(stmt.group.stmt)
      if payload:
        return payload
    elif stmt.HasField('alt'):
      for choice in stmt.alt.choice:
        payload = return_payload(choice.stmt)
        if payload:
          return payload
    elif stmt.HasField('ret'):
      return stmt.ret.payload
    else:
      raise Exception('No statement!', stmt.WhichOneof('stmt'))
