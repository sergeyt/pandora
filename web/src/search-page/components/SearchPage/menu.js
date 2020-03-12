import DashboardIcon from '@material-ui/icons/Dashboard';
import ShoppingCartIcon from '@material-ui/icons/ShoppingCart';
import PeopleIcon from '@material-ui/icons/People';
import BarChartIcon from '@material-ui/icons/BarChart';
import LayersIcon from '@material-ui/icons/Layers';
import AssignmentIcon from '@material-ui/icons/Assignment';

export const primaryItems = [
    {
        icon: DashboardIcon,
        text: "Search Documents",
        link: "/",
    },
    {
        icon: ShoppingCartIcon,
        text: "Upload Files",
        link: "",
    },
    {
        icon: PeopleIcon,
        text: "Customers",
        link: "/Dashboard/mini",
    },
    {
        icon: BarChartIcon,
        text: "Reports",
        link: "/Dashboard/mini",
    },
    {
        icon: LayersIcon,
        text: "Integrations",
        link: "/Dashboard/mini",
    },
];

export const secondaryItems = [
    {
        icon: AssignmentIcon,
        text: "Current month",
        link: "/Dashboard/mini",
    },
    {
        icon: AssignmentIcon,
        text: "Last quarter",
        link: "/Dashboard/mini",
    },
    {
        icon: AssignmentIcon,
        text: "Year-end sale",
        link: "/Dashboard/mini",
    },
];

const categories = [
    {
        items: primaryItems,
    },
    {
        title: "",
        items: secondaryItems,
    },
];

export default categories;
