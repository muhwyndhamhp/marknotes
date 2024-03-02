import {common, createLowlight} from 'lowlight'

import js from 'highlight.js/lib/languages/javascript'
import xml from 'highlight.js/lib/languages/xml'
import go from 'highlight.js/lib/languages/go'
import css from 'highlight.js/lib/languages/css'
import python from 'highlight.js/lib/languages/python'
import java from 'highlight.js/lib/languages/java'
import kotlin from 'highlight.js/lib/languages/kotlin'
import lua from 'highlight.js/lib/languages/lua'
import ts from 'highlight.js/lib/languages/typescript'
import c from 'highlight.js/lib/languages/c'
import rust from 'highlight.js/lib/languages/rust'
import php from 'highlight.js/lib/languages/php'
import ruby from 'highlight.js/lib/languages/ruby'
import swift from 'highlight.js/lib/languages/swift'
import dart from 'highlight.js/lib/languages/dart'
import yaml from 'highlight.js/lib/languages/yaml'
import json from 'highlight.js/lib/languages/json'
import shell from 'highlight.js/lib/languages/shell'
import sql from 'highlight.js/lib/languages/sql'


export const lowlight = createLowlight(common)

lowlight.register({js})
lowlight.register({xml})
lowlight.register({go})
lowlight.register({css})
lowlight.register({python})
lowlight.register({java})
lowlight.register({kotlin})
lowlight.register({lua})
lowlight.register({ts})
lowlight.register({c})
lowlight.register({rust})
lowlight.register({php})
lowlight.register({ruby})
lowlight.register({swift})
lowlight.register({dart})
lowlight.register({yaml})
lowlight.register({json})
lowlight.register({shell})
lowlight.register({sql})
