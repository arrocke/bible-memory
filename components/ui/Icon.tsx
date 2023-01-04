import { ComponentProps } from 'react'
const { library, config } = require('@fortawesome/fontawesome-svg-core')
import { faLeftLong, faPlus, faMinus, faPencil, faCheck, faClose } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import '@fortawesome/fontawesome-svg-core/styles.css'
config.autoAddCss = false

library.add(faLeftLong, faPlus, faMinus, faPencil, faCheck, faClose)

export type IconProps = ComponentProps<typeof FontAwesomeIcon>

export default function Icon(props: IconProps) {
  return <FontAwesomeIcon {...props} />
}