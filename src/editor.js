import { Color } from '@tiptap/extension-color'
import { Editor } from '@tiptap/core'
import ListItem from '@tiptap/extension-list-item'
import StarterKit from '@tiptap/starter-kit'
import TextStyle from '@tiptap/extension-text-style'
import CodeBlockLowlight from '@tiptap/extension-code-block-lowlight'
import Image from '@tiptap/extension-image'
import Link from '@tiptap/extension-link'
import Placeholder from '@tiptap/extension-placeholder'
import Mention from '@tiptap/extension-mention'
import tippy from 'tippy.js'

const lowlight = require('./unexported/lowlight.js').lowlight;

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
      CodeBlockLowlight.configure({ lowlight }),
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
      Mention.configure({
         suggestion: {
            char: "#",
            render: () => {
               let popup

               return {
                  onStart: (props) => {
                     popup = tippy(document.body, {
                        getReferenceClientRect: props.clientRect,
                        appendTo: () => document.body,
                        content: `
                         <div class="dropdown dropdown-open">
                           <ul class="p-4 shadow menu dropdown-content z-[1] bg-base-100 rounded-box w-64">
                              <li class="text-lg text-base-content">
                                 Start Write to Search...
                              </li>
                           </ul>
                        </div>
                        `,
                        allowHTML: true,
                        showOnCreate: true,
                        trigger: 'manual',
                        placement: 'bottom-start',
                     })
                  },
                  onUpdate: () => {
                     content = editor.state
                        .selection
                        .$head
                        .parent
                        .textContent
                        .replace("#", "")

                     debounceUpdate(content, popup);
                  },
                  onExit() {
                     popup.destroy()
                  }
               }
            }
         }
      }),
   ],
   editorProps: {
      attributes: {
         class: 'caret-rose-800 w-full mx-auto prose prose-slate lg:prose-xl md:prose-lg dark:prose-invert prose-pre:bg-slate-900 prose-pre:w-full prose-pre:text-white focus:outline-none prose-h2:text-primary prose-h3:text-secondary prose-h4:text-accent prose-em:text-secondary prose-strong:text-primary prose-strong:font-extrabold prose-a:font-extrabold prose-a:text-accent',
      },
   },
})

function debounce(callback, delay) {
   let timerID;
   return (...args) => {
      clearTimeout(timerID);
      timerID = setTimeout(() => {
         callback(...args);
      }, delay);
   }
}

const debounceUpdate = debounce((content, popup) => {
   fetch(`/dashboard/tags?tag=${content}`, {
      method: 'GET',
   }).then((response) => {
      return response.text()
   }).then((data) => {
      popup.setProps({
         content: data,
      })
   })
}, 500)

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

