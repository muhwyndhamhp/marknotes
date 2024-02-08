function backPress() {
   history.back();
};

document.onload = function () {
   checkTheme()
}

//TODO: Does not work 
document.addEventListener('checkTheme', function (evt) {
   checkTheme()
})

function checkTheme() {
   if (localStorage.theme === 'dark'
      || (!('theme' in localStorage)
         && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
      document.documentElement.classList.add('dark')
   } else {
      document.documentElement.classList.remove('dark')
   }
}

function allowDrop(ev) {
   ev.preventDefault()
}

function upload(ev, url) {
   ev.preventDefault()
   file = ev.dataTransfer.files[0]

   Swal.showLoading()

   const formData = new FormData()
   formData.append("file", file)

   fetch(url, {
      method: "POST",
      body: formData,
      contentType: "multipart/form-data"
   }).then((response) => {
      return response.text()
   }).then(afterUpload);
}

function afterUpload(rawData) {
   data = JSON.parse(rawData)
   elm = document.getElementById("form-content")
   imgMark = `![img](${data.data.url})`
   if (elm.selectionStart || elm.selectionStart == '0') {
      var startPos = elm.selectionStart;
      var endPos = elm.selectionEnd;

      elm.value = elm.value.substring(0, startPos)
         + imgMark
         + elm.value.substring(endPos, elm.value.length);
   }
   else {
      elm.value += imgMark
   }

   Swal.close()
}

const content = `
<h2>
  Hi there,
</h2>
<p>
  this is a <em>basic</em> example of <strong>tiptap</strong>. Sure, there are all kind of basic text styles you’d probably expect from a text editor. But wait until you see the lists:
</p>
<ul>
  <li>
    That’s a bullet list with one …
  </li>
  <li>
    … or two list items.
  </li>
</ul>
<p>
  Isn’t that great? And all of that is editable. But wait, there’s more. Let’s try a code block:
</p>
<pre><code class="language-javascript">
for (var i=1; i <= 20; i++)
{
  if (i % 15 == 0)
    console.log("FizzBuzz");
  else if (i % 3 == 0)
    console.log("Fizz");
  else if (i % 5 == 0)
    console.log("Buzz");
  else
    console.log(i);
}
</code></pre>
<p>
  I know, I know, this is impressive. It’s only the tip of the iceberg though. Give it a try and click a little bit around. Don’t forget to check the other examples too.
</p>
<blockquote>
  Wow, that’s amazing. Good work, boy! 👏
  <br />
  — Mom
</blockquote>
`

import { Color } from '@tiptap/extension-color'
import { Editor } from '@tiptap/core'
import ListItem from '@tiptap/extension-list-item'
import StarterKit from '@tiptap/starter-kit'
import TextStyle from '@tiptap/extension-text-style'
import CodeBlockLowlight from '@tiptap/extension-code-block-lowlight'
import {common, createLowlight} from 'lowlight'
import js from 'highlight.js/lib/languages/javascript'

const lowlight = createLowlight(common)

lowlight.register({js})

const editor = new Editor({
   element: document.querySelector('#code-editor'),
   extensions: [
      Color.configure({ types: [TextStyle.name, ListItem.name] }),
      TextStyle.configure({ types: [ListItem.name] }),
      StarterKit.configure({
         bulletList: {
            keepMarks: true,
            keepAttributes: false, // TODO : Making this as `false` becase marks are not preserved when I try to preserve attrs, awaiting a bit of help
         },
         orderedList: {
            keepMarks: true,
            keepAttributes: false, // TODO : Making this as `false` becase marks are not preserved when I try to preserve attrs, awaiting a bit of help
         },
      }),
      CodeBlockLowlight.configure({ lowlight }),
   ],
   editorProps: {
      attributes: {
         class: 'prose prose-slate lg:prose-xl md:prose-lg dark:prose-invert prose-h2:bg-clip-text prose-h2:text-transparent prose-h2:bg-gradient-to-r prose-h2:from-pink-600 prose-h2:via-purple-600 prose-h2:to-indigo-600 prose-h3:bg-clip-text  prose-h3:text-transparent  prose-h3:bg-gradient-to-r  prose-h3:from-teal-600  prose-h3:via-sky-600  prose-h3:to-indigo-600 prose-pre:bg-slate-900 prose-pre:w-full prose-pre:text-white',
      },
   },
   content: content,
})

