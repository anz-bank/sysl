import os


class FileWriter(object):
  def __init__(self, out_dir, package, entities):
    self.out_dir = out_dir
    self.package = package
    self.entities = entities

  def __call__(self, w, out_path):
    out_path = os.path.join(self.out_dir, out_path)
    try:
      os.makedirs(os.path.dirname(out_path))
    except:
      pass
    open(out_path, 'w').write(str(w))

  def java(self, w, name, package):
    assert name
    if self.entities:
      assert name in self.entities, ('Unexpected entity generated: ' + name +
        ' (check BUILD has sysl_model(..., entities=[..., "' + name + ', ...], ...))')
      self.entities.remove(name)
    self(w, os.path.join(self.package.replace('.', '/'), name + u'.java'))
