+++
# AUTOGENERATED BY byexample/generate.go
title= "Integration Diagrams"
draft= false
description= ""
layout= "byexample"
weight = 6
topic = "Diagrams"
Images = [
  
  "/assets/byexample/images/integration-diagrams5.png",
  
]

ID = "integration-diagrams"
Segs = [[
  
      {CodeEmpty= true,CodeLeading= true,CodeRun= false,CodeRendered="""""",DocsRendered= """<p>In this example will use a simple system and start using the sysl command to generate diagrams.</p>
""", CodeForJs = """"""},

      {CodeEmpty= false,CodeLeading= true,CodeRun= true,CodeRendered="""<div class="highlight"><pre><span class="nx">IntegratedSystem</span><span class="p">:</span>
    <span class="nx">integrated_endpoint_1</span><span class="p">:</span>
        <span class="nx">System1</span> <span class="o">&lt;-</span> <span class="nx">endpoint</span>
    <span class="nx">integrated_endpoint_2</span><span class="p">:</span>
        <span class="nx">System2</span> <span class="o">&lt;-</span> <span class="nx">endpoint</span>
</pre></div>
""",DocsRendered= """""", CodeForJs = """IntegratedSystem:
    integrated_endpoint_1:
        System1 <- endpoint
    integrated_endpoint_2:
        System2 <- endpoint
"""},

      {CodeEmpty= false,CodeLeading= true,CodeRun= false,CodeRendered="""<div class="highlight"><pre><span class="nx">System1</span><span class="p">:</span>
    <span class="nx">endpoint</span><span class="p">:</span> <span class="o">...</span>
<span class="nx">System2</span><span class="p">:</span>
    <span class="nx">endpoint</span><span class="p">:</span> <span class="o">...</span>
</pre></div>
""",DocsRendered= """""", CodeForJs = """System1:
    endpoint: ...
System2:
    endpoint: ...
"""},

      {CodeEmpty= false,CodeLeading= false,CodeRun= false,CodeRendered="""<div class="highlight"><pre><span class="nx">Project</span> <span class="p">[</span><span class="nx">appfmt</span><span class="p">=</span><span class="s">&quot;%(appname)&quot;</span><span class="p">]:</span>
    <span class="nx">_</span><span class="p">:</span>
        <span class="nx">IntegratedSystem</span>
        <span class="nx">System1</span>
        <span class="nx">System2</span>
</pre></div>
""",DocsRendered= """""", CodeForJs = """Project [appfmt="%(appname)"]:
    _:
        IntegratedSystem
        System1
        System2
"""},

],
[
  
      {CodeEmpty= false,CodeLeading= true,CodeRun= false,CodeRendered="""<div class="highlight"><pre><span class="go">export SYSL_PLANTUML=http://www.plantuml.com/plantuml</span>
</pre></div>
""",DocsRendered= """<p>First, make sure to set the environment variable SYSL_PLANTUML</p>
""", CodeForJs = """"""},

      {CodeEmpty= false,CodeLeading= true,CodeRun= true,CodeRendered="""<div class="highlight"><pre><span class="go">sysl ints -o project.png --project Project project.sysl</span>
</pre></div>
""",DocsRendered= """<p>Now run the sysl sd (sequence diagram) command</p>
""", CodeForJs = """"""},

      {CodeEmpty= true,CodeLeading= true,CodeRun= false,CodeRendered="""""",DocsRendered= """<p><code>-o</code> is the output file</p>
""", CodeForJs = """"""},

      {CodeEmpty= true,CodeLeading= true,CodeRun= false,CodeRendered="""""",DocsRendered= """<p><code>-s</code> specifies a starting endpoint for the sequence diagram to initiate</p>
""", CodeForJs = """"""},

      {CodeEmpty= true,CodeLeading= true,CodeRun= false,CodeRendered="""""",DocsRendered= """<p><code>project.sysl</code> is the input sysl file</p>
""", CodeForJs = """"""},

      {CodeEmpty= true,CodeLeading= false,CodeRun= false,CodeRendered="""""",DocsRendered= """<p>project.png:</p>
""", CodeForJs = """"""},

],

]
+++

