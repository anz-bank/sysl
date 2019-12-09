+++
title= "hello-world3"
draft= false
description= ""
layout= "byexample"
Images = [
  
]

ID = "hello-world3.md"
Segs = [[
  
      {CodeEmpty= false,CodeLeading= true,CodeRun= false,CodeRendered="""<div class="highlight"><pre><span class="nx">HelloWorld</span><span class="p">:</span>
    <span class="o">/</span><span class="nx">endpoint1</span><span class="p">:</span>
        <span class="nx">GET</span><span class="p">:</span>
</pre></div>
""",DocsRendered= """<p>Our first program will print the classic &ldquo;hello world&rdquo;message. Here&rsquo;s the full source code. this is more text for doing more stuff</p>
"""},

      {CodeEmpty= false,CodeLeading= true,CodeRun= false,CodeRendered="""<div class="highlight"><pre>            <span class="k">return</span> <span class="kt">string</span>
</pre></div>
""",DocsRendered= """<p>Here we can specify a return type</p>
"""},

      {CodeEmpty= false,CodeLeading= false,CodeRun= false,CodeRendered="""<div class="highlight"><pre><span class="nx">proj</span><span class="p">:</span>
    <span class="nx">_</span><span class="p">:</span>
        <span class="nx">HelloWorld</span>
</pre></div>
""",DocsRendered= """<p>message. Here&rsquo;s the full source code. this is more text for doing more stuff</p>
"""},

],
[
  
      {CodeEmpty= false,CodeLeading= true,CodeRun= false,CodeRendered="""<div class="highlight"><pre><span class="gp">$</span> go run hello-world.go
<span class="go">hello world</span>
</pre></div>
""",DocsRendered= """<p>To run the program, put the code in <code>hello-world.go</code> anduse  go run.</p>
"""},

      {CodeEmpty= false,CodeLeading= true,CodeRun= false,CodeRendered="""<div class="highlight"><pre><span class="gp">$</span> go build hello-world.go
<span class="gp">$</span> ls
<span class="go">hello-world    hello-world.go</span>
</pre></div>
""",DocsRendered= """<p>Sometimes we&rsquo;ll want to build our programs intobinaries. We can do this using <code>go build</code>.</p>
"""},

      {CodeEmpty= false,CodeLeading= true,CodeRun= false,CodeRendered="""<div class="highlight"><pre><span class="gp">$</span> ./hello-world
<span class="go">hello world</span>
</pre></div>
""",DocsRendered= """<p>We can then execute the built binary directly.</p>
"""},

      {CodeEmpty= true,CodeLeading= false,CodeRun= false,CodeRendered="""""",DocsRendered= """<p>Now that we can run and build basic Go programs, let&rsquo;slearn more about the language.</p>
"""},

],

]
+++


