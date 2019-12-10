+++
# AUTOGENERATED BY byexample/generate.go
title= "Data-types"
draft= false
description= ""
layout= "byexample"
weight = 3
topic = "Basics"
Images = [
  
]

ID = "data-types"
Segs = [[
  
      {CodeEmpty= true,CodeLeading= true,CodeRun= false,CodeRendered="""""",DocsRendered= """<p>Our first program will make a simple &ldquo;Hello world&rdquo; application with two endpoints</p>
""", CodeForJs = """"""},

      {CodeEmpty= false,CodeLeading= true,CodeRun= true,CodeRendered="""<div class="highlight"><pre><span class="nx">HelloWorld</span><span class="p">:</span>
</pre></div>
""",DocsRendered= """<p>Specify an application called <code>HelloWorld</code></p>
""", CodeForJs = """HelloWorld:
"""},

      {CodeEmpty= false,CodeLeading= true,CodeRun= false,CodeRendered="""<div class="highlight"><pre>    <span class="k">!type</span> <span class="nx">Request</span><span class="p">:</span>
        <span class="nx">userId</span> <span class="p">&lt;:</span> <span class="kt">int</span>
        <span class="nx">id</span> <span class="p">&lt;:</span> <span class="kt">int</span>
        <span class="nx">title</span> <span class="p">&lt;:</span> <span class="kt">string</span>
        <span class="nx">completed</span> <span class="p">&lt;:</span> <span class="kt">bool</span>
    
    <span class="k">!type</span> <span class="nx">ErrorResponse</span><span class="p">:</span>
        <span class="nx">status</span> <span class="p">&lt;:</span> <span class="kt">string</span>
</pre></div>
""",DocsRendered= """<p>Specify composite types with &ldquo;!type&rdquo; followed by type fields</p>
""", CodeForJs = """    !type Request:
        userId <: int
        id <: int
        title <: string
        completed <: bool
    
    !type ErrorResponse:
        status <: string
"""},

      {CodeEmpty= false,CodeLeading= true,CodeRun= false,CodeRendered="""<div class="highlight"><pre>    <span class="k">!alias</span> <span class="nx">Requests</span><span class="p">:</span>
        <span class="nx">sequence</span> <span class="nx">of</span> <span class="nx">Request</span>
</pre></div>
""",DocsRendered= """<p>Use the <code>!alias</code> keyword to alias to another name</p>
""", CodeForJs = """    !alias Requests:
        sequence of Request
"""},

      {CodeEmpty= false,CodeLeading= true,CodeRun= false,CodeRendered="""<div class="highlight"><pre>    <span class="nx">endpoint</span><span class="p">(</span><span class="nx">input</span> <span class="p">&lt;:</span> <span class="nx">Request</span><span class="p">):</span>
</pre></div>
""",DocsRendered= """<p>Specify an endpoint as the next indent.</p>
""", CodeForJs = """    endpoint(input <: Request):
"""},

      {CodeEmpty= false,CodeLeading= true,CodeRun= false,CodeRendered="""<div class="highlight"><pre>        <span class="k">return</span> <span class="nx">Requests</span>
</pre></div>
""",DocsRendered= """<p>Specify a return type for the endpoint. Anything after the return is considered a payload.</p>
""", CodeForJs = """        return Requests
"""},

      {CodeEmpty= true,CodeLeading= false,CodeRun= false,CodeRendered="""""",DocsRendered= """<p>Seeing that we have only the simplest sysl files and no interactions between services we cannot run any meaningful commands.</p>
""", CodeForJs = """"""},

],

]
+++

