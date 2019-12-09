+++
title= "todos"
draft= false
description= ""
layout= "byexample"


ID = "todos.md"
Segs = [[
  
      {CodeEmpty= false,CodeLeading= true,CodeRun= false,CodeRendered="""<div class="highlight"><pre><span class="nx">Todos</span><span class="p">:</span>
</pre></div>
""",DocsRendered= """<p>Here we define an Application</p>
"""},

      {CodeEmpty= false,CodeLeading= true,CodeRun= false,CodeRendered="""<div class="highlight"><pre>  <span class="p">!</span><span class="kd">type</span> <span class="nx">Todo</span><span class="p">:</span>
    <span class="nx">userId</span> <span class="p">&lt;:</span> <span class="kt">int</span>
    <span class="nx">id</span> <span class="p">&lt;:</span> <span class="kt">int</span>
    <span class="nx">title</span> <span class="p">&lt;:</span> <span class="kt">string</span>
    <span class="nx">completed</span> <span class="p">&lt;:</span> <span class="kt">bool</span>
</pre></div>
""",DocsRendered= """<p>Here we can define a type with different fields</p>
"""},

      {CodeEmpty= false,CodeLeading= true,CodeRun= false,CodeRendered="""<div class="highlight"><pre>  <span class="p">!</span><span class="kd">type</span> <span class="nx">Post</span><span class="p">:</span>
    <span class="nx">userId</span> <span class="p">&lt;:</span> <span class="kt">int</span>
    <span class="nx">id</span> <span class="p">&lt;:</span> <span class="kt">int</span>
    <span class="nx">title</span> <span class="p">&lt;:</span> <span class="kt">string</span>
    <span class="nx">body</span> <span class="p">&lt;:</span> <span class="kt">string</span>
</pre></div>
""",DocsRendered= """<p>Use the <code>!alias</code> keyword to alias to another name or a set or sequence</p>
"""},

      {CodeEmpty= false,CodeLeading= true,CodeRun= false,CodeRendered="""<div class="highlight"><pre>  <span class="p">!</span><span class="nx">alias</span> <span class="nx">Posts</span><span class="p">:</span>
    <span class="nx">sequence</span> <span class="nx">of</span> <span class="nx">Post</span>
</pre></div>
""",DocsRendered= """"""},

      {CodeEmpty= false,CodeLeading= true,CodeRun= false,CodeRendered="""<div class="highlight"><pre>  <span class="p">!</span><span class="kd">type</span> <span class="nx">ErrorResponse</span><span class="p">:</span>
    <span class="nx">status</span> <span class="p">&lt;:</span> <span class="kt">string</span>
</pre></div>
""",DocsRendered= """"""},

      {CodeEmpty= false,CodeLeading= true,CodeRun= false,CodeRendered="""<div class="highlight"><pre>  <span class="p">!</span><span class="kd">type</span> <span class="nx">ResourceNotFoundError</span><span class="p">:</span>
    <span class="nx">status</span> <span class="p">&lt;:</span> <span class="kt">string</span>
</pre></div>
""",DocsRendered= """"""},

      {CodeEmpty= false,CodeLeading= true,CodeRun= false,CodeRendered="""<div class="highlight"><pre>  <span class="o">/</span><span class="nx">todos</span><span class="p">:</span>
    <span class="o">/</span><span class="p">{</span><span class="nx">id</span><span class="p">&lt;:</span><span class="kt">int</span><span class="p">}:</span>
      <span class="nx">GET</span><span class="p">:</span>
        <span class="k">if</span> <span class="nx">notfound</span><span class="p">:</span>
          <span class="k">return</span> <span class="mi">404</span> <span class="p">&lt;:</span> <span class="nx">ResourceNotFoundError</span>
        <span class="k">else</span> <span class="k">if</span> <span class="nx">failed</span><span class="p">:</span>
          <span class="k">return</span> <span class="mi">500</span> <span class="p">&lt;:</span> <span class="nx">ErrorResponse</span>
        <span class="k">else</span><span class="p">:</span>    
          <span class="k">return</span> <span class="mi">200</span> <span class="p">&lt;:</span> <span class="nx">Todo</span> 
</pre></div>
""",DocsRendered= """<p>Here we define the todos endpoint with a get reponse</p>
"""},

      {CodeEmpty= false,CodeLeading= true,CodeRun= false,CodeRendered="""<div class="highlight"><pre>  <span class="o">/</span><span class="nx">posts</span><span class="p">:</span>
    <span class="nx">GET</span><span class="p">:</span>
      <span class="k">if</span> <span class="nx">notfound</span><span class="p">:</span>
        <span class="k">return</span> <span class="mi">404</span> <span class="p">&lt;:</span> <span class="nx">ResourceNotFoundError</span>
      <span class="k">else</span> <span class="k">if</span> <span class="nx">failed</span><span class="p">:</span>
        <span class="k">return</span> <span class="mi">500</span> <span class="p">&lt;:</span> <span class="nx">ErrorResponse</span>
      <span class="k">else</span><span class="p">:</span>    
        <span class="k">return</span> <span class="mi">200</span> <span class="p">&lt;:</span> <span class="nx">Posts</span>
</pre></div>
""",DocsRendered= """"""},

      {CodeEmpty= false,CodeLeading= false,CodeRun= false,CodeRendered="""<div class="highlight"><pre>  <span class="o">/</span><span class="nx">comments</span><span class="p">:</span>
    <span class="nx">GET</span> <span class="err">?</span><span class="nx">postId</span><span class="p">=</span><span class="kt">int</span><span class="p">:</span>
      <span class="k">return</span> <span class="nx">Posts</span>
      
    <span class="nx">POST</span> <span class="p">(</span><span class="nx">newPost</span> <span class="p">&lt;:</span> <span class="nx">Post</span> <span class="p">[</span><span class="err">~</span><span class="nx">body</span><span class="p">]):</span>
      <span class="k">return</span> <span class="nx">Post</span>
</pre></div>
""",DocsRendered= """"""},

],

]
+++


