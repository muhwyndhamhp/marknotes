import { Color } from '@tiptap/extension-color'
import { Editor } from '@tiptap/core'
import ListItem from '@tiptap/extension-list-item'
import StarterKit from '@tiptap/starter-kit'
import TextStyle from '@tiptap/extension-text-style'
import CodeBlockLowlight from '@tiptap/extension-code-block-lowlight'
import Image from '@tiptap/extension-image'
import Link from '@tiptap/extension-link'
import Placeholder from '@tiptap/extension-placeholder'
import Youtube from '@tiptap/extension-youtube'

const lowlight = require('./unexported/lowlight.js').lowlight;
const hashtag = require('./unexported/hashtag.js').HashTag;

const editorProps = require('./unexported/editor_props.js').EditorProps;

export const editor = new Editor({
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
      CodeBlockLowlight.configure({ 
         lowlight,
         HTMLAttributes: {
            class: 'mockup-code rounded-badge',
         }
      }),
      Image.configure({
         inline: true,
         HTMLAttributes: {
            class: 'max-h-96 mx-auto'
         }
      }), 
      Link.configure({
        openOnClick: true,
        autolink: true,
        linkOnPaste: true,
      }),
      Placeholder.configure({
         considerAnyAsEmpty: true,
         placeholder: 'Write your thought here...',
      }),
      Youtube.configure({
         inline: false,
         width: 512,
         height: 288,
         modestBranding: 'true',
         progressBarColor: 'white',
         HTMLAttributes: {
            class: 'mx-auto'
         }
      }),
      hashtag,
   ],
   editorProps: editorProps,
   onCreate: () => {
      const encodedContent = document.getElementById('code-editor').children[0].innerHTML
      document.getElementById('content').value = encodedContent
   },

   onUpdate: ({ editor }) => {
      const encodedContent = document.getElementById('code-editor').children[0].innerHTML
      console.log(encodedContent)
      document.getElementById('content').value = encodedContent

      const tags = editor.getJSON().content
         .filter((node) => node.type === 'paragraph')
         .filter((node) => node.content !== undefined)
         .filter((node) => node.content.length > 0)
         .flatMap((node) => 
            node
            .content
            .filter((child) => child !== undefined)
            .filter((child) => child.type === 'mention')
            .map((child) => child.attrs.id)
         )      

      const uqTags = [...new Set(tags)]
      document.getElementById('tags').value = uqTags.join(',')
   },
   content: window.content,
   
})


window.editor = editor


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

