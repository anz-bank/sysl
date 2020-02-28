require 'rubygems/package'
require 'zlib'

TAR_LONGLINK = '././@LongLink'
tar_gz_archive = File.join(ENV['HOME'], 'Library', 'Caches', 'Homebrew', 'sysl--0.6.3.tar.gz')
destination = File.join(ENV['HOME'], 'Library', 'Caches', 'Homebrew')

Gem::Package::TarReader.new( Zlib::GzipReader.open tar_gz_archive ) do |tar|
  dest = nil
  tar.each do |entry|
    if entry.full_name == TAR_LONGLINK
      dest = File.join destination, entry.read.strip
      next
    end
    dest ||= File.join destination, entry.full_name
    if entry.directory?
      File.delete dest if File.file? dest
      FileUtils.mkdir_p dest, :mode => entry.header.mode, :verbose => false
    elsif entry.file?
      FileUtils.rm_rf dest if File.directory? dest
      File.open dest, "wb" do |f|
        f.print entry.read
      end
      FileUtils.chmod entry.header.mode, dest, :verbose => false
    elsif entry.header.typeflag == '2' 
      File.symlink entry.header.linkname, dest
    end
    dest = nil
  end
end

Src = File.join(ENV['HOME'], 'Library', 'Caches', 'Homebrew', 'sysl-0.6.3', 'cmd', 'sysl')
Dir.chdir Src
system("go", "build")
Dest = File.join(ENV['HOME'],'go', 'bin')
FileUtils.cp(File.join(Src, 'sysl') ,Dest)
