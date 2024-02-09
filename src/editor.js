import { Color } from '@tiptap/extension-color'
import { Editor } from '@tiptap/core'
import ListItem from '@tiptap/extension-list-item'
import StarterKit from '@tiptap/starter-kit'
import TextStyle from '@tiptap/extension-text-style'
import CodeBlockLowlight from '@tiptap/extension-code-block-lowlight'
import Image from '@tiptap/extension-image'

const content = require('./unexported/content_placeholder.js').content;
const lowlight = require('./unexported/lowlight.js').lowlight;

const editor = new Editor({
   element: document.querySelector('#code-editor'),
   extensions: [
      Color.configure({ types: [TextStyle.name, ListItem.name] }),
      TextStyle.configure({ types: [ListItem.name] }),
      StarterKit.configure({
         bulletList: {
            keepMarks: true,
            keepAttributes: false, 
         },
         orderedList: {
            keepMarks: true,
            keepAttributes: false, 
         },
      }),
      CodeBlockLowlight.configure({ lowlight }),
      Image.configure({
         inline: true,
         HTMLAttributes: {
            class: 'max-w-96 mx-auto'
         }
      })
   ],
   editorProps: {
      attributes: {
         class: 'prose prose-slate lg:prose-xl md:prose-lg dark:prose-invert prose-h2:bg-clip-text prose-h2:text-transparent prose-h2:bg-gradient-to-r prose-h2:from-pink-600 prose-h2:via-purple-600 prose-h2:to-indigo-600 prose-h3:bg-clip-text  prose-h3:text-transparent  prose-h3:bg-gradient-to-r  prose-h3:from-teal-600  prose-h3:via-sky-600  prose-h3:to-indigo-600 prose-pre:bg-slate-900 prose-pre:w-full prose-pre:text-white',
      },
   },
   content: content,
})


window.allowDrop = function (ev) {
   ev.preventDefault()
}

window.upload = function(ev, url) {
   ev.preventDefault()


   if(ev.dataTransfer.files.length === 0) {
      return
   }

   file = ev.dataTransfer.files[0]

   console.log(file)

   Swal.showLoading()

   const formData = new FormData()
   formData.append("file", file)

   fetch(url, {
      method: "POST",
      body: formData,
      contentType: "multipart/form-data"
   }).then((response) => {
      return response.text()
   }).then(window.afterUpload);
}

window.afterUpload = function(rawData) {
   data = JSON.parse(rawData)
   editor.chain().focus().setImage({ src: data.data.url }).run()
   Swal.close()
}

